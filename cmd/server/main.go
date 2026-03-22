package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	ConfigRuntime()
	InitDB()
	StartGin()
}

// ConfigRuntime sets the number of operating system threads.
func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

// InitDB initializes the database connection.
func InitDB() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Create tables
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test (
		id SERIAL PRIMARY KEY,
		data TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected and table created")
}



// StartGin starts gin web server with setting router.
func StartGin() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	// Session store
	store := cookie.NewStore([]byte("secret-key"))
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("resources/*.html")
	router.GET("/", index)
	router.GET("/form", AuthRequired(), formGET)
	router.POST("/submit", AuthRequired(), formPOST)
	router.GET("/db", dbVersion)
	router.GET("/login", loginGET)
	router.POST("/login", loginPOST)
	router.GET("/logout", logoutGET)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
