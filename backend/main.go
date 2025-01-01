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
	ID                    int64     `json:"id"`
	Role                  string    `json:"role"`
	OtherRole             string    `json:"otherRole"`
	CmsUsage              string    `json:"cmsUsage"`
	OtherCmsUsage         string    `json:"otherCmsUsage"`
	Features              Features  `json:"features"`
	BetaInterest          bool      `json:"betaInterest"`
	Email                 string    `json:"email"`
	CreatedAt             time.Time `json:"created_at"`
	BiggestFrustrations   string    `json:"biggestFrustrations"`
	SpecificProblems      string    `json:"specificProblems"`
	UsageFrequency        string    `json:"usageFrequency"`
	PrimaryPurpose        string    `json:"primaryPurpose"`
	Platforms             string    `json:"platforms"`
	CmsPreference         string    `json:"cmsPreference"`
	WishedFeatures        string    `json:"wishedFeatures"`
	WorkflowImportance    string    `json:"workflowImportance"`
	TeamSize              string    `json:"teamSize"`
	CollaborationFrequency string   `json:"collaborationFrequency"`
	PricingSensitivity    string    `json:"pricingSensitivity"`
	PricingModel          string    `json:"pricingModel"`
	Integrations          string    `json:"integrations"`
	IntegrationImportance string    `json:"integrationImportance"`
	ContentTypes          string    `json:"contentTypes"`
	CustomFormats         string    `json:"customFormats"`
	FeedbackSuggestions   string    `json:"feedbackSuggestions"`
	ExcitementFactors     string    `json:"excitementFactors"`
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

	// Create or update tables
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

	// Update existing table schema to add new fields
	alterTableQueries := []string{
		"ALTER TABLE survey_responses ADD COLUMN biggest_frustrations TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN specific_problems TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN usage_frequency TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN primary_purpose TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN platforms TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN cms_preference TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN wished_features TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN workflow_importance TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN team_size TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN collaboration_frequency TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN pricing_sensitivity TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN pricing_model TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN integrations TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN integration_importance TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN content_types TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN custom_formats TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN feedback_suggestions TEXT;",
		"ALTER TABLE survey_responses ADD COLUMN excitement_factors TEXT;",
	}

	for _, query := range alterTableQueries {
		_, err = db.Exec(query)
		if err != nil {
			// Log but do not fail if the column already exists
			log.Printf("Error updating table schema: %v", err)
		}
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
            beta_interest, email, biggest_frustrations, specific_problems,
            usage_frequency, primary_purpose, platforms, cms_preference,
            wished_features, workflow_importance, team_size, collaboration_frequency,
            pricing_sensitivity, pricing_model, integrations, integration_importance,
            content_types, custom_formats, feedback_suggestions, excitement_factors
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		response.Role, response.OtherRole, response.CmsUsage, response.OtherCmsUsage,
		response.Features.Offline, response.Features.Collaboration, response.Features.AssetManagement,
		response.Features.PdfHandling, response.Features.VersionControl, response.Features.Workflows,
		response.BetaInterest, response.Email, response.BiggestFrustrations, response.SpecificProblems,
		response.UsageFrequency, response.PrimaryPurpose, response.Platforms, response.CmsPreference,
		response.WishedFeatures, response.WorkflowImportance, response.TeamSize, response.CollaborationFrequency,
		response.PricingSensitivity, response.PricingModel, response.Integrations, response.IntegrationImportance,
		response.ContentTypes, response.CustomFormats, response.FeedbackSuggestions, response.ExcitementFactors,
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
		beta_interest, email, biggest_frustrations, specific_problems,
		usage_frequency, primary_purpose, platforms, cms_preference,
		wished_features, workflow_importance, team_size, collaboration_frequency,
		pricing_sensitivity, pricing_model, integrations, integration_importance,
		content_types, custom_formats, feedback_suggestions, excitement_factors,
		created_at FROM survey_responses`)
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
			&response.BetaInterest, &response.Email, &response.BiggestFrustrations,
			&response.SpecificProblems, &response.UsageFrequency, &response.PrimaryPurpose,
			&response.Platforms, &response.CmsPreference, &response.WishedFeatures,
			&response.WorkflowImportance, &response.TeamSize, &response.CollaborationFrequency,
			&response.PricingSensitivity, &response.PricingModel, &response.Integrations,
			&response.IntegrationImportance, &response.ContentTypes, &response.CustomFormats,
			&response.FeedbackSuggestions, &response.ExcitementFactors, &response.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, responses)
}

func deleteResult(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM survey_responses WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Result deleted"})
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

func verifyToken(c *gin.Context) {
	// The AuthMiddleware already verified the token
	// We just need to return a success response
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":    true,
		"username": username,
	})
}

func getTrustedProxies() []string {
	// Get from environment variable, fallback to Docker network ranges
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		return strings.Split(trustedProxies, ",")
	}

	// Default Docker network ranges
	return []string{
		"172.16.0.0/12",  // Docker default
		"192.168.0.0/16", // Docker default
		"10.0.0.0/8",     // Docker default
		"127.0.0.1",      // Localhost
	}
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

	// Set trusted proxies
	if err := router.SetTrustedProxies(getTrustedProxies()); err != nil {
		log.Printf("Warning: Failed to set trusted proxies: %v", err)
	}

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
		authorized.GET("/verify", verifyToken)
		authorized.DELETE("/results/:id", deleteResult)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090" // Change default port to 8090
	}

	router.Run(":" + port)
}
