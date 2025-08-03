package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestEcommerceAPI(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "E-commerce API Suite")
}

var _ = ginkgo.Describe("E-commerce API", func() {
	var router *gin.Engine
	var testDB *gorm.DB

	ginkgo.BeforeEach(func() {
		// Setup test database
		var err error
		testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Set global db to test db
		db = testDB

		// Migrate schema
		db.AutoMigrate(&User{}, &Item{}, &Cart{}, &CartItem{}, &Order{}, &OrderItem{})

		// Seed test data
		seedTestData()

		// Setup router
		gin.SetMode(gin.TestMode)
		router = gin.New()
		setupRoutes(router)
	})

	ginkgo.Describe("User Management", func() {
		ginkgo.It("should create a new user", func() {
			payload := map[string]string{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			}
			jsonPayload, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			gomega.Expect(w.Code).To(gomega.Equal(http.StatusCreated))
		})

		ginkgo.It("should login user and return token", func() {
			// First create a user
			user := User{
				Username: "logintest",
				Email:    "login@example.com",
				Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			}
			db.Create(&user)

			payload := map[string]string{
				"username": "logintest",
				"password": "password",
			}
			jsonPayload, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			gomega.Expect(w.Code).To(gomega.Equal(http.StatusOK))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			gomega.Expect(response["token"]).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("Cart Management", func() {
		var token string

		ginkgo.BeforeEach(func() {
			// Create and login user
			user := User{
				Username: "cartuser",
				Email:    "cart@example.com",
				Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
			}
			db.Create(&user)
			token = generateToken(user.ID)
			db.Model(&user).Update("token", token)
		})

		ginkgo.It("should add item to cart", func() {
			payload := map[string]interface{}{
				"item_id":  1,
				"quantity": 2,
			}
			jsonPayload, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			gomega.Expect(w.Code).To(gomega.Equal(http.StatusOK))
		})

		ginkgo.It("should get user cart", func() {
			req, _ := http.NewRequest("GET", "/carts", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Should return 404 if no cart exists yet
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusNotFound))
		})
	})

	ginkgo.Describe("Order Management", func() {
		var token string
		var userID uint

		ginkgo.BeforeEach(func() {
			// Create and login user
			user := User{
				Username: "orderuser",
				Email:    "order@example.com",
				Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
			}
			db.Create(&user)
			userID = user.ID
			token = generateToken(user.ID)
			db.Model(&user).Update("token", token)

			// Create cart with items
			cart := Cart{UserID: userID}
			db.Create(&cart)
			cartItem := CartItem{
				CartID:   cart.ID,
				ItemID:   1,
				Quantity: 2,
			}
			db.Create(&cartItem)
		})

		ginkgo.It("should create order from cart", func() {
			req, _ := http.NewRequest("POST", "/orders", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			gomega.Expect(w.Code).To(gomega.Equal(http.StatusCreated))
		})

		ginkgo.It("should list user orders", func() {
			req, _ := http.NewRequest("GET", "/orders", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			gomega.Expect(w.Code).To(gomega.Equal(http.StatusOK))
		})
	})
})

func seedTestData() {
	items := []Item{
		{Name: "Test Laptop", Description: "Test laptop", Price: 999.99},
		{Name: "Test Mouse", Description: "Test mouse", Price: 29.99},
	}

	for _, item := range items {
		db.Create(&item)
	}
}
