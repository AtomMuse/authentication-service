package main

import (
	_ "atommuse/backend/authentication-service/cmd/authentication/doc"
	handlers "atommuse/backend/authentication-service/handler"
	"atommuse/backend/authentication-service/pkg/repository"
	services "atommuse/backend/authentication-service/pkg/service"
	"fmt"
	"log"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title			Authentication Service API
// @version		v0
// @description	Authentication Service สำหรับขอจัดการเกี่ยวกับ Authentication
// @schemes		http
func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	fmt.Print(mongoURI)
	mongoDBConn, err := repository.NewMongoDBConnection(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoDBConn.Close(); err != nil {
			log.Fatalf("Error closing MongoDB connection: %v", err)
		}
	}()

	// Get a database instance
	dbName := "atommuse"
	db := mongoDBConn.GetDatabase(dbName)

	// Repositories
	userRepo := repository.NewUserRepository(db)

	// Services
	authService := services.NewAuthService(userRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	// Swagger documentation route
	router := gin.Default()
	router.Use(cors.New(config))

	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/register", authHandler.Register)
	}

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
