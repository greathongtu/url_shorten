package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/exp/rand"
)

type URLMapping struct {
	ID        int    `json:"id"`
	LongURL   string `json:"long_url"`
	ShortCode string `json:"short_code"`
}

var db *sql.DB

func main() {
	// Initialize database connection
	initDB()

	// Set up Gin router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Add your frontend URL here
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Define routes
	r.POST("/shorten", shortenURL)
	r.GET("/:shortCode", redirectURL)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func initDB() {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@db:5432/urlshortener?sslmode=disable"
	}

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS url_mappings (
			id SERIAL PRIMARY KEY,
			long_url TEXT NOT NULL,
			short_code VARCHAR(10) UNIQUE NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func shortenURL(c *gin.Context) {
	var input struct {
		LongURL string `json:"long_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var shortCode string
	var err error
	for attempts := 0; attempts < 5; attempts++ {
		shortCode = generateShortCode()
		_, err = db.Exec("INSERT INTO url_mappings (long_url, short_code) VALUES ($1, $2)", input.LongURL, shortCode)
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "duplicate key value") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store URL"})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate unique short code"})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"short_url": fmt.Sprintf("http://%s/%s", c.Request.Host, shortCode)})
	c.JSON(http.StatusOK, gin.H{"short_code": shortCode})
}

func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCode := make([]byte, 6)
	for i := range shortCode {
		shortCode[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortCode)
}

func redirectURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var longURL string
	err := db.QueryRow("SELECT long_url FROM url_mappings WHERE short_code = $1", shortCode).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		}
		return
	}

	c.Redirect(http.StatusMovedPermanently, longURL)
}
