package main

import (
	_ "atommuse/backend/authentication-service/cmd/authentication/doc"
	handlers "atommuse/backend/authentication-service/handler"
	"atommuse/backend/authentication-service/pkg/model"
	"atommuse/backend/authentication-service/pkg/repository"
	services "atommuse/backend/authentication-service/pkg/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title						Authentication Service API
// @version					v0
// @description				Authentication Service สำหรับขอจัดการเกี่ยวกับ Authentication
// @schemes					http
//
// @SecurityDefinitions.apikey	BearerAuth
// @In							header
// @Name						Authorization
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
	userService := services.NewUserService(userRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	// Swagger documentation route
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // Replace "*" with allowed origins
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/register", authHandler.Register)
	}
	apiRoutes := router.Group("/api")
	{
		apiRoutes.PUT("/user/:id", authMiddleware("exhibitor"), userHandler.UpdateUserByID)
		apiRoutes.PUT("/user/change-password", authMiddleware("exhibitor"), userHandler.ChangePassword)
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

// authMiddleware is middleware to validate the token and check the role
func authMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		secretKey := os.Getenv("SECRET_KEY") // Corrected env variable name
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		// Parse the token
		claims := &model.JwtCustomClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			// Check the token signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the secret key for validation
			return []byte(secretKey), nil
		})

		fmt.Println(parsedToken.Valid)

		// Handle token parsing errors
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			fmt.Println("Token parsing error:", err)
			return
		}

		// Check if the token is valid
		if !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			fmt.Println("Invalid token")
			return
		}

		// Set user ID in context
		c.Set("user_id", claims.UserID)

		// Check if the role matches
		if claims.Role != "admin" && claims.Role != "exhibitor" && claims.Role != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			fmt.Println("Insufficient permissions")
			return
		}

		// Continue down the chain to handler etc
		c.Next()
	}
}
