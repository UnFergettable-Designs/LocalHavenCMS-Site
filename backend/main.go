package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/time/rate"
)

type SurveyResponse struct {
	ID            int64     `json:"id"`
	Role          string    `json:"role"`
	OtherRole     string    `json:"otherRole"`
	CmsUsage      string    `json:"cmsUsage"`
	OtherCmsUsage string    `json:"otherCmsUsage"`
	Features      Features  `json:"features"`
	BetaInterest  bool      `json:"betaInterest"`
	Email         string    `json:"email"`
	CreatedAt     time.Time `json:"created_at"`
}

type Features struct {
	Offline         int `json:"offline"`
	Collaboration   int `json:"collaboration"`
	AssetManagement int `json:"assetManagement"`
	PdfHandling     int `json:"pdfHandling"`
	VersionControl  int `json:"versionControl"`
	Workflows       int `json:"workflows"`
}

var db *sql.DB

// Add input validation
var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	limiter    = rate.NewLimiter(rate.Every(time.Second), 10) // 10 requests per second
	clients    = make(map[string]*rate.Limiter)
	mu         sync.Mutex
)

func init() {
	// Only try to load .env file in development
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: No .env file not found, using environment variables")
		}
	}
}

func getEnvWithFallback(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func initDB() {
	var err error
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "/app/data/localhavencms.db"
	}

	// Create data directory with full permissions
	dataDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dataDir, 0777); err != nil {
		log.Printf("Error creating data directory: %v", err)
		// Try to create in current directory as fallback
		dbPath = "localhavencms.db"
	}

	log.Printf("Opening database at: %s", dbPath)

	// Try to create an empty file first to test permissions
	if _, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		log.Fatalf("Cannot create database file: %v", err)
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Test database connection
	if err = db.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Printf("Successfully connected to database at: %s", dbPath)

	// Create tables
	createTables := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS survey_responses (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        role TEXT NOT NULL,
        other_role TEXT,
        cms_usage TEXT NOT NULL,
        other_cms_usage TEXT,
        offline_rating INTEGER,
        collaboration_rating INTEGER,
        asset_management_rating INTEGER,
        pdf_handling_rating INTEGER,
        version_control_rating INTEGER,
        workflows_rating INTEGER,
        beta_interest BOOLEAN,
        email TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatal(err)
	}
}

func validateSurveyResponse(r *SurveyResponse) error {
	if r.Role == "" {
		return fmt.Errorf("role is required")
	}
	if r.CmsUsage == "" {
		return fmt.Errorf("CMS usage is required")
	}
	if r.BetaInterest && (r.Email == "" || !emailRegex.MatchString(r.Email)) {
		return fmt.Errorf("valid email is required for beta program")
	}
	return nil
}

func getClientLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if limiter, exists := clients[ip]; exists {
		return limiter
	}

	newLimiter := rate.NewLimiter(rate.Every(time.Minute), 60) // 60 requests per minute per IP
	clients[ip] = newLimiter
	return newLimiter
}

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getClientLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func submitSurvey(c *gin.Context) {
	var response SurveyResponse
	if err := c.BindJSON(&response); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	if err := validateSurveyResponse(&response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received survey response: %+v", response)

	stmt, err := db.Prepare(`
        INSERT INTO survey_responses(
            role, other_role, cms_usage, other_cms_usage,
            offline_rating, collaboration_rating, asset_management_rating,
            pdf_handling_rating, version_control_rating, workflows_rating,
            beta_interest, email
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		response.Role, response.OtherRole, response.CmsUsage, response.OtherCmsUsage,
		response.Features.Offline, response.Features.Collaboration, response.Features.AssetManagement,
		response.Features.PdfHandling, response.Features.VersionControl, response.Features.Workflows,
		response.BetaInterest, response.Email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Survey response recorded"})
}

func getSurveyResults(c *gin.Context) {
	rows, err := db.Query(`SELECT id, role, other_role, cms_usage, other_cms_usage,
		offline_rating, collaboration_rating, asset_management_rating,
		pdf_handling_rating, version_control_rating, workflows_rating,
		beta_interest, email, created_at FROM survey_responses`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var responses []SurveyResponse
	for rows.Next() {
		var response SurveyResponse
		err := rows.Scan(&response.ID, &response.Role, &response.OtherRole,
			&response.CmsUsage, &response.OtherCmsUsage,
			&response.Features.Offline, &response.Features.Collaboration,
			&response.Features.AssetManagement, &response.Features.PdfHandling,
			&response.Features.VersionControl, &response.Features.Workflows,
			&response.BetaInterest, &response.Email, &response.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, responses)
}

func login(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	if user.Username != os.Getenv("ADMIN_USERNAME") || user.Password != os.Getenv("ADMIN_PASSWORD") {
		time.Sleep(time.Second) // Prevent timing attacks
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Remove 'Bearer ' prefix if present
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Add claims to context if needed
			c.Set("username", claims["username"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// Set Gin mode based on environment
	if getEnvWithFallback("ENVIRONMENT", "development") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Verify required environment variables
	requiredEnvVars := []string{"JWT_SECRET", "ADMIN_USERNAME", "ADMIN_PASSWORD"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Required environment variable %s is not set", envVar)
		}
	}

	initDB()
	defer db.Close()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		allowedOrigins := []string{
			"http://localhost:3000",
			"https://localhavencms.com",
			"https://www.localhavencms.com",
			"https://api.localhavencms.com",
		}

		// Log the incoming request details for debugging
		log.Printf("Incoming request from: %s to: %s", c.Request.Host, c.Request.URL.Path)

		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Set security headers
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		c.Next()
	})

	// Use rate limiting
	router.Use(rateLimitMiddleware())

	// Public routes
	router.POST("/survey", submitSurvey)
	router.POST("/login", login)

	// Add explicit health check logging
	router.GET("/health", func(c *gin.Context) {
		log.Printf("Health check from: %s", c.Request.Host)
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"host":   c.Request.Host,
		})
	})

	// Protected routes
	authorized := router.Group("/")
	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/results", getSurveyResults)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
