package main

import (
	"errors"
	"github.com/ichtrojan/go-todo/middleware"
	"github.com/ichtrojan/go-todo/models"
	"github.com/ichtrojan/go-todo/routes"
	"github.com/ichtrojan/thoth"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	logger, _ := thoth.Init("log")

	if err := godotenv.Load(); err != nil {
		logger.Log(errors.New("no .env file found"))
		log.Fatal("No .env file found")
	}

	// Initialize session store with secure key from environment variable
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		sessionKey = "todo-app-secret-key-change-in-production"
		logger.Log(errors.New("SESSION_KEY not set in .env, using default key"))
		log.Println("Warning: SESSION_KEY not set in .env, using default key")
	}
	
	// Set the session store with the session key
	middleware.Store = sessions.NewCookieStore([]byte(sessionKey))

	// Initialize database tables
	models.InitUserTable() // Create users table if it doesn't exist
	models.InitTodoTable() // Create todos table with user_id foreign key

	port, exist := os.LookupEnv("PORT")

	if !exist {
		logger.Log(errors.New("PORT not set in .env"))
		log.Fatal("PORT not set in .env")
	}

	// Start the server
	log.Println("Server starting on port " + port)
	err := http.ListenAndServe(":"+port, routes.Init())

	if err != nil {
		logger.Log(errors.New("server error: " + err.Error()))
		log.Fatal(err)
	}
}