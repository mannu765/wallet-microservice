package main

import (
    "log"
	"os"
    "wallet-microservice/internal/database"
    "wallet-microservice/internal/handlers"
    "wallet-microservice/internal/repositories"
    "wallet-microservice/internal/services"
    
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    // Connect to database
    database.Connect()
    database.Migrate()
    
    // Initialize layers
    walletRepo := repositories.NewWalletRepository()
    walletService := services.NewWalletService(walletRepo)
    walletHandler := handlers.NewWalletHandler(walletService)
    
    // Setup Gin router
    router := gin.Default()
    
    // Add middleware
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    
    // Add CORS middleware
    router.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })
    
    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "OK",
            "service": "wallet-microservice",
        })
    })
    
    // Register routes
    walletHandler.RegisterRoutes(router)
    
    // Start server
    port := getEnv("PORT", "8080")
    log.Printf("Server starting on port %s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
