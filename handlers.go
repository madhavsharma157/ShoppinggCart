package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// User handlers
func createUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func loginUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token (simple implementation - in production use JWT)
	token := generateToken(user.ID)
	
	// Update user token (single device login)
	db.Model(&user).Update("token", token)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}

func listUsers(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Cart handlers
func addToCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		ItemID   uint `json:"item_id" binding:"required"`
		Quantity int  `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if item exists
	var item Item
	if err := db.First(&item, req.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Get or create cart for user
	var cart Cart
	if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		// Create new cart
		cart = Cart{UserID: userID}
		db.Create(&cart)
	}

	// Check if item already in cart
	var cartItem CartItem
	if err := db.Where("cart_id = ? AND item_id = ?", cart.ID, req.ItemID).First(&cartItem).Error; err != nil {
		// Add new item to cart
		cartItem = CartItem{
			CartID:   cart.ID,
			ItemID:   req.ItemID,
			Quantity: req.Quantity,
		}
		db.Create(&cartItem)
	} else {
		// Update quantity
		cartItem.Quantity += req.Quantity
		db.Save(&cartItem)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item added to cart",
		"cart_item": cartItem,
	})
}

func getCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var cart Cart
	if err := db.Preload("Items.Item").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func listCarts(c *gin.Context) {
	var carts []Cart
	db.Preload("User").Preload("Items.Item").Find(&carts)
	c.JSON(http.StatusOK, gin.H{"carts": carts})
}

// Order handlers
func createOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Get user's cart
	var cart Cart
	if err := db.Preload("Items.Item").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Calculate total
	var total float64
	for _, cartItem := range cart.Items {
		total += cartItem.Item.Price * float64(cartItem.Quantity)
	}

	// Create order
	order := Order{
		UserID: userID,
		Total:  total,
		Status: "pending",
	}
	db.Create(&order)

	// Create order items
	for _, cartItem := range cart.Items {
		orderItem := OrderItem{
			OrderID:  order.ID,
			ItemID:   cartItem.ItemID,
			Quantity: cartItem.Quantity,
			Price:    cartItem.Item.Price,
		}
		db.Create(&orderItem)
	}

	// Clear cart
	db.Where("cart_id = ?", cart.ID).Delete(&CartItem{})
	db.Delete(&cart)

	// Load order with items
	db.Preload("Items.Item").First(&order, order.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})
}

func listOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var orders []Order
	db.Preload("Items.Item").Where("user_id = ?", userID).Find(&orders)
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// Item handlers
func listItems(c *gin.Context) {
	var items []Item
	
	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	db.Offset(offset).Limit(limit).Find(&items)
	
	var total int64
	db.Model(&Item{}).Count(&total)
	
	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
