package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	var err error
	db, err = gorm.Open(sqlite.Open("ecommerce.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&User{}, &Item{}, &Cart{}, &CartItem{}, &Order{}, &OrderItem{})

	// Seed some initial items
	seedItems()

	// Initialize Gin router
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Routes
	setupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}

func setupRoutes(r *gin.Engine) {
	// User routes
	r.POST("/users", createUser)
	r.POST("/users/login", loginUser)
	r.GET("/users", listUsers)

	// Protected routes (require authentication)
	protected := r.Group("/")
	protected.Use(authMiddleware())
	{
		// Cart routes
		protected.POST("/carts", addToCart)
		protected.GET("/carts", getCart)
		protected.GET("/carts/all", listCarts)

		// Order routes
		protected.POST("/orders", createOrder)
		protected.GET("/orders", listOrders)

		// Items routes
		protected.GET("/items", listItems)
	}
}

func seedItems() {
	// Check if items already exist
	var count int64
	db.Model(&Item{}).Count(&count)
	if count > 0 {
		return
	}

	items := []Item{
		{Name: "Laptop", Description: "High-performance laptop", Price: 999.99},
		{Name: "Mouse", Description: "Wireless mouse", Price: 29.99},
		{Name: "Keyboard", Description: "Mechanical keyboard", Price: 79.99},
		{Name: "Monitor", Description: "4K monitor", Price: 299.99},
		{Name: "Headphones", Description: "Noise-canceling headphones", Price: 199.99},
	}

	for _, item := range items {
		db.Create(&item)
	}
}
