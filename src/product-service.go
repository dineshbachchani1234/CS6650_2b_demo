package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// Product matches the Product schema in api.yaml. All six fields are required.
// The `binding` tags tell Gin how to validate incoming JSON.
type Product struct {
	ProductID    int    `json:"product_id"    binding:"required,min=1"`
	SKU          string `json:"sku"           binding:"required,min=1,max=100"`
	Manufacturer string `json:"manufacturer"  binding:"required,min=1,max=200"`
	CategoryID   int    `json:"category_id"   binding:"required,min=1"`
	Weight       int    `json:"weight"        binding:"min=0"`
	SomeOtherID  int    `json:"some_other_id" binding:"required,min=1"`
}

// Error matches the Error schema in api.yaml.
type Error struct {
	ErrorCode string `json:"error"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
}

// In-memory product store, guarded by a RWMutex so concurrent requests are safe.
var (
	products = make(map[int]Product)
	mu       sync.RWMutex
)

func main() {
	router := gin.Default()

	router.GET("/products/:productId", getProduct)
	router.POST("/products/:productId/details", addProductDetails)

	router.Run(":8080")
}

// GET /products/{productId}
// 200 = found, 404 = not found, 400 = bad id (not a valid positive int)
func getProduct(c *gin.Context) {
	idStr := c.Param("productId")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, Error{
			ErrorCode: "INVALID_INPUT",
			Message:   "Product ID must be a positive integer",
			Details:   "Received: " + idStr,
		})
		return
	}

	mu.RLock()
	product, exists := products[id]
	mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, Error{
			ErrorCode: "NOT_FOUND",
			Message:   "Product not found",
			Details:   "No product with ID " + idStr,
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// POST /products/{productId}/details
// 204 = success, 400 = invalid input, 404 = product id mismatch handling
func addProductDetails(c *gin.Context) {
	idStr := c.Param("productId")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, Error{
			ErrorCode: "INVALID_INPUT",
			Message:   "Product ID must be a positive integer",
			Details:   "Received: " + idStr,
		})
		return
	}

	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, Error{
			ErrorCode: "INVALID_INPUT",
			Message:   "The provided input data is invalid",
			Details:   err.Error(),
		})
		return
	}

	// Store (or update) the product's details under the path id.
	mu.Lock()
	products[id] = product
	mu.Unlock()

	// 204 No Content on success, as the spec requires.
	c.Status(http.StatusNoContent)
}
