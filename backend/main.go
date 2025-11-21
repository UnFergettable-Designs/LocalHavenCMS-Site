package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/time/rate"
)

type SurveyResponse struct {
	ID                            string    `json:"id"`
	Role                          string    `json:"role"`
	OtherRole                     string    `json:"otherRole,omitempty"`
	CmsUsage                      string    `json:"cmsUsage"`
	OtherCmsUsage                 string    `json:"otherCmsUsage,omitempty"`
	Features                      Features  `json:"features"`
	BetaInterest                  bool      `json:"betaInterest"`
	Email                         string    `json:"email,omitempty"`
	CreatedAt                     time.Time `json:"createdAt"`
	BiggestFrustrations           string    `json:"biggestFrustrations"`
	SpecificProblems              string    `json:"specificProblems"`
	UsageFrequency                string    `json:"usageFrequency"`
	PrimaryPurpose                string    `json:"primaryPurpose"`
	Platforms                     string    `json:"platforms"`
	CmsPreference                 string    `json:"cmsPreference"`
	WishedFeatures                string    `json:"wishedFeatures"`
	WorkflowImportance            string    `json:"workflowImportance"`
	TeamSize                      string    `json:"teamSize"`
	CollaborationFrequency        string    `json:"collaborationFrequency"`
	PricingSensitivity            string    `json:"pricingSensitivity"`
	PricingModel                  string    `json:"pricingModel"`
	Integrations                  string    `json:"integrations"`
	IntegrationImportance         string    `json:"integrationImportance"`
	ContentTypes                  string    `json:"contentTypes"`
	CustomFormats                 string    `json:"customFormats"`
	FeedbackSuggestions           string    `json:"feedbackSuggestions"`
	ExcitementFactors             string    `json:"excitementFactors"`
	CollaborationChallenges       string    `json:"collaborationChallenges"`
	OfflineWorkFrequency          string    `json:"offlineWorkFrequency"`
	OfflineWorkarounds            string    `json:"offlineWorkarounds"`
	CurrentChangeConflictHandling string    `json:"currentChangeConflictHandling"`
	VersionControlChallenges      string    `json:"versionControlChallenges"`
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

// Add caching for survey results
var (
	resultsCache   []SurveyResponse
	cacheTimestamp time.Time
	cacheDuration  = 5 * time.Minute
	cacheMutex     sync.RWMutex
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

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "data/localhavencms.db")
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Check if table exists
	var tableExists bool
	err = db.QueryRow(`SELECT COUNT(*) FROM sqlite_master 
		WHERE type='table' AND name='survey_responses'`).Scan(&tableExists)

	if err != nil || !tableExists {
		// Create new table if it doesn't exist
		return createInitialTable()
	}

	// Check if we need to migrate
	var columnCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('survey_responses') 
		WHERE name IN ('biggest_frustrations', 'specific_problems')`).Scan(&columnCount)

	if err != nil || columnCount < 2 {
		return migrateTable()
	}

	return nil
}

func createInitialTable() error {
	_, err := db.Exec(`CREATE TABLE survey_responses (
		id TEXT PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		role TEXT NOT NULL,
		other_role TEXT,
		cms_usage TEXT NOT NULL,
		other_cms_usage TEXT,
		-- Features as individual columns
		offline INTEGER,
		collaboration INTEGER,
		asset_management INTEGER,
		pdf_handling INTEGER,
		version_control INTEGER,
		workflows INTEGER,
		beta_interest BOOLEAN NOT NULL,
		email TEXT,
		biggest_frustrations TEXT,
		specific_problems TEXT,
		usage_frequency TEXT,
		primary_purpose TEXT,
		platforms TEXT,
		cms_preference TEXT,
		wished_features TEXT,
		workflow_importance TEXT,
		team_size TEXT,
		collaboration_frequency TEXT,
		pricing_sensitivity TEXT,
		pricing_model TEXT,
		integrations TEXT,
		integration_importance TEXT,
		content_types TEXT,
		custom_formats TEXT,
		feedback_suggestions TEXT,
		excitement_factors TEXT,
		collaboration_challenges TEXT,
		offline_work_frequency TEXT,
		offline_workarounds TEXT,
		current_change_conflict_handling TEXT,
		version_control_challenges TEXT
	)`)
	return err
}

func migrateTable() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create new table with all fields
	_, err = tx.Exec(`CREATE TABLE survey_responses_new (
        id TEXT PRIMARY KEY,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        role TEXT NOT NULL,
        other_role TEXT,
        cms_usage TEXT NOT NULL,
        other_cms_usage TEXT,
        features JSON,
        beta_interest BOOLEAN NOT NULL,
        email TEXT,
        biggest_frustrations TEXT,
        specific_problems TEXT,
        usage_frequency TEXT,
        primary_purpose TEXT,
        platforms TEXT,
        cms_preference TEXT,
        wished_features TEXT,
        workflow_importance TEXT,
        team_size TEXT,
        collaboration_frequency TEXT,
        pricing_sensitivity TEXT,
        pricing_model TEXT,
        integrations TEXT,
        integration_importance TEXT,
        content_types TEXT,
        custom_formats TEXT,
        feedback_suggestions TEXT,
        excitement_factors TEXT,
        collaboration_challenges TEXT,
        offline_work_frequency TEXT,
        offline_workarounds TEXT,
        current_change_conflict_handling TEXT,
        version_control_challenges TEXT
    )`)
	if err != nil {
		return err
	}

	// Copy data including all feature fields
	_, err = tx.Exec(`
        INSERT INTO survey_responses_new (
            id, created_at, role, other_role, cms_usage, other_cms_usage,
            offline, collaboration, asset_management,
            pdf_handling, version_control, workflows,
            beta_interest, email
        )
        SELECT 
            id, created_at, role, other_role, cms_usage, other_cms_usage,
            COALESCE(offline, 0), 
            COALESCE(collaboration, 0),
            COALESCE(asset_management, 0),
            COALESCE(pdf_handling, 0),
            COALESCE(version_control, 0),
            COALESCE(workflows, 0),
            beta_interest, email
        FROM survey_responses
    `)

	// Complete migration
	if err != nil {
		return err
	}

	return tx.Commit()
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

func endpointRateLimiter(limit rate.Limit, burst int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)
	var mu sync.Mutex

	return func(c *gin.Context) {
		ip := c.ClientIP()
		endpoint := c.FullPath()
		key := fmt.Sprintf("%s:%s", ip, endpoint)

		mu.Lock()
		limiter, exists := limiters[key]
		if !exists {
			limiter = rate.NewLimiter(limit, burst)
			limiters[key] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func submitSurvey(c *gin.Context) {
	var survey SurveyResponse
	if err := c.BindJSON(&survey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	survey.ID = uuid.New().String()
	survey.CreatedAt = time.Now()

	// Log the survey data
	log.Printf("Submitting survey: %+v\n", survey)
	log.Printf("Features: %+v\n", survey.Features)

	stmt, err := db.Prepare(`
		INSERT INTO survey_responses (
			id, created_at, role, other_role, cms_usage, other_cms_usage,
			offline, collaboration, asset_management,
			pdf_handling, version_control, workflows,
			beta_interest, email,
			biggest_frustrations, specific_problems,
			usage_frequency, primary_purpose, platforms,
			cms_preference, wished_features, workflow_importance,
			team_size, collaboration_frequency,
			pricing_sensitivity, pricing_model,
			integrations, integration_importance,
			content_types, custom_formats,
			feedback_suggestions, excitement_factors,
			collaboration_challenges, offline_work_frequency, offline_workarounds,
			current_change_conflict_handling, version_control_challenges
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?,
			?, ?, ?,
			?, ?,
			?, ?,
			?, ?, ?,
			?, ?, ?,
			?, ?,
			?, ?,
			?, ?,
			?, ?,
			 ?, ?, ?, ?
		)
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		survey.ID, survey.CreatedAt, survey.Role, survey.OtherRole,
		survey.CmsUsage, survey.OtherCmsUsage,
		survey.Features.Offline, survey.Features.Collaboration,
		survey.Features.AssetManagement, survey.Features.PdfHandling,
		survey.Features.VersionControl, survey.Features.Workflows,
		survey.BetaInterest, survey.Email,
		survey.BiggestFrustrations, survey.SpecificProblems,
		survey.UsageFrequency, survey.PrimaryPurpose,
		survey.Platforms, survey.CmsPreference,
		survey.WishedFeatures, survey.WorkflowImportance,
		survey.TeamSize, survey.CollaborationFrequency,
		survey.PricingSensitivity, survey.PricingModel,
		survey.Integrations, survey.IntegrationImportance,
		survey.ContentTypes, survey.CustomFormats,
		survey.FeedbackSuggestions, survey.ExcitementFactors,
		survey.CollaborationChallenges,
		survey.OfflineWorkFrequency,
		survey.OfflineWorkarounds,
		survey.CurrentChangeConflictHandling,
		survey.VersionControlChallenges,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, survey)
}

func getSurveyResults(c *gin.Context) {
	cacheMutex.RLock()
	if time.Since(cacheTimestamp) < cacheDuration && resultsCache != nil {
		defer cacheMutex.RUnlock()
		c.JSON(http.StatusOK, resultsCache)
		return
	}
	cacheMutex.RUnlock()

	rows, err := db.Query(`
		SELECT id, role, other_role, cms_usage, other_cms_usage,
		offline, collaboration, asset_management,
		pdf_handling, version_control, workflows,
		beta_interest, email, biggest_frustrations, specific_problems,
		usage_frequency, primary_purpose, platforms, cms_preference,
		wished_features, workflow_importance, team_size, collaboration_frequency,
		pricing_sensitivity, pricing_model, integrations, integration_importance,
		content_types, custom_formats, feedback_suggestions, excitement_factors,
		collaboration_challenges, offline_work_frequency, offline_workarounds,
		current_change_conflict_handling, version_control_challenges,
		created_at 
		FROM survey_responses`)
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
			&response.FeedbackSuggestions, &response.ExcitementFactors, &response.CollaborationChallenges,
			&response.OfflineWorkFrequency, &response.OfflineWorkarounds, &response.CurrentChangeConflictHandling,
			&response.VersionControlChallenges, &response.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		responses = append(responses, response)
	}

	cacheMutex.Lock()
	resultsCache = responses
	cacheTimestamp = time.Now()
	cacheMutex.Unlock()

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

func getTrustedProxies() ([]string, error) {
	// Get from environment variable, fallback to Docker network ranges
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		proxies := strings.Split(trustedProxies, ",")
		// Validate each proxy address/range
		for _, proxy := range proxies {
			proxy = strings.TrimSpace(proxy)
			if !strings.Contains(proxy, "/") && net.ParseIP(proxy) == nil {
				return nil, fmt.Errorf("invalid proxy address: %s", proxy)
			}
		}
		return proxies, nil
	}

	// Default secure Docker network ranges
	return []string{
		"172.16.0.0/12",  // Docker default bridge network
		"192.168.0.0/16", // Docker user-defined networks
		"10.0.0.0/8",     // Docker overlay networks
		"127.0.0.1",      // Localhost
	}, nil
}

// Add metrics collection
type Metrics struct {
	TotalResponses       int                `json:"totalResponses"`
	BetaInterestCount    int                `json:"betaInterestCount"`
	AverageFeatureScores map[string]float64 `json:"averageFeatureScores"`
	UsageFrequencyStats  map[string]int     `json:"usageFrequencyStats"`
	TeamSizeDistribution map[string]int     `json:"teamSizeDistribution"`
	PricingPreferences   map[string]int     `json:"pricingPreferences"`
}

func getMetrics(c *gin.Context) {
	metrics := Metrics{
		AverageFeatureScores: make(map[string]float64),
		UsageFrequencyStats:  make(map[string]int),
		TeamSizeDistribution: make(map[string]int),
		PricingPreferences:   make(map[string]int),
	}

	var (
		offlineScore, collaborationScore, assetScore float64
		pdfScore, vcScore, workflowScore             float64
	)

	// Get basic counts and feature averages
	err := db.QueryRow(`
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN beta_interest = 1 THEN 1 ELSE 0 END) as beta_count,
			AVG(offline) as avg_offline,
			AVG(collaboration) as avg_collab,
			AVG(asset_management) as avg_asset,
			AVG(pdf_handling) as avg_pdf,
			AVG(version_control) as avg_vc,
			AVG(workflows) as avg_workflow
		FROM survey_responses
	`).Scan(
		&metrics.TotalResponses,
		&metrics.BetaInterestCount,
		&offlineScore,
		&collaborationScore,
		&assetScore,
		&pdfScore,
		&vcScore,
		&workflowScore,
	)

	// Then assign to map
	metrics.AverageFeatureScores["offline"] = offlineScore
	metrics.AverageFeatureScores["collaboration"] = collaborationScore
	metrics.AverageFeatureScores["assetManagement"] = assetScore
	metrics.AverageFeatureScores["pdfHandling"] = pdfScore
	metrics.AverageFeatureScores["versionControl"] = vcScore
	metrics.AverageFeatureScores["workflows"] = workflowScore
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get usage frequency distribution
	rows, err := db.Query(`
		SELECT usage_frequency, COUNT(*) as count 
		FROM survey_responses 
		GROUP BY usage_frequency
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var frequency string
		var count int
		if err := rows.Scan(&frequency, &count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		metrics.UsageFrequencyStats[frequency] = count
	}

	// Get team size distribution
	rows, err = db.Query(`
		SELECT team_size, COUNT(*) as count 
		FROM survey_responses 
		GROUP BY team_size
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var size string
		var count int
		if err := rows.Scan(&size, &count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		metrics.TeamSizeDistribution[size] = count
	}

	// Get pricing preferences
	rows, err = db.Query(`
		SELECT pricing_model, COUNT(*) as count 
		FROM survey_responses 
		GROUP BY pricing_model
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var model string
		var count int
		if err := rows.Scan(&model, &count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		metrics.PricingPreferences[model] = count
	}

	// Round feature scores to 2 decimal places
	for key, value := range metrics.AverageFeatureScores {
		metrics.AverageFeatureScores[key] = float64(int(value*100)) / 100
	}

	c.JSON(http.StatusOK, metrics)
}

func setupRouter(env string) *gin.Engine {
	r := gin.Default()

	// Disable rate limiting for preview environment
	if os.Getenv("RATE_LIMIT_DISABLED") == "true" {
		r.GET("/health", func(c *gin.Context) {
			c.Header("X-Environment", "preview")
			c.String(200, "healthy")
		})
	} else {
		// Use rate limiting
		r.Use(rateLimitMiddleware())

		// Public routes
		r.POST("/survey", endpointRateLimiter(rate.Every(time.Minute), 5), submitSurvey)
		r.POST("/login", endpointRateLimiter(rate.Every(time.Minute), 3), login)

		// Add explicit health check logging
		r.GET("/health", func(c *gin.Context) {
			log.Printf("Health check from: %s", c.Request.Host)
			c.JSON(http.StatusOK, gin.H{
				"status": "healthy",
				"host":   c.Request.Host,
			})
		})

		// Protected routes
		authorized := r.Group("/")
		authorized.Use(AuthMiddleware())
		{
			authorized.GET("/results", getSurveyResults)
			authorized.GET("/verify", verifyToken)
			authorized.DELETE("/results/:id", deleteResult)
			authorized.GET("/metrics", getMetrics)
		}
	}

	return r
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

	// Set trusted proxies with proper error handling
	trustedProxies, err := getTrustedProxies()
	if err != nil {
		log.Fatalf("Failed to configure trusted proxies: %v", err)
	}

	router := setupRouter(os.Getenv("ENVIRONMENT"))

	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Log the configured trusted proxies
	log.Printf("Configured trusted proxies: %v", trustedProxies)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090" // Change default port to 8090
	}

	router.Run(":" + port)
}
