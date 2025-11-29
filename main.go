package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Harshitttttttt/go-todo-jwt/auth"
	"github.com/Harshitttttttt/go-todo-jwt/db"
	"github.com/Harshitttttttt/go-todo-jwt/handlers"
	"github.com/Harshitttttttt/go-todo-jwt/middleware"
	"github.com/Harshitttttttt/go-todo-jwt/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// loadEnv loads the env variables from the .env
func loadEnv() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using env variables")
	}

	// Check required variables
	requiredVars := []string{"DATABASE_URL", "JWT_SECRET"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Required env variable %s is not set", v)
		}
	}
}

func main() {
	// Load env variables
	loadEnv()

	// Connect to the database
	database, err := db.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	r := mux.NewRouter()

	// Create Repositories
	userRepo := models.NewUserRepository(database)
	refreshTokenRepo := models.NewRefreshTokenRepository(database)

	// Create Services
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, os.Getenv("JWT_SECRET"), 15*time.Minute)

	// Create Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)

	// Public Routes
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/api/auth/refresh", authHandler.RefreshToken).Methods("POST")

	// Protected Routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(authService))

	protected.HandleFunc("/profile", userHandler.Profile).Methods("GET")

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
