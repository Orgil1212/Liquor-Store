package controllers

import (
	"liquor-store/config"
	"liquor-store/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 🛒 Сагсанд бүтээгдэхүүн нэмэх
func AddToCart(c *gin.Context) {
	var cartItem models.Cart

	// 📌 JSON хүсэлтийг зөв уншиж байгаа эсэхийг шалгах
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 📌 Бүтээгдэхүүнийг сагсанд нэмэх
	if err := config.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to cart"})
		return
	}

	// ✅ JSON форматтай хариу буцаах
	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart"})
}

// 🛒 Хэрэглэгчийн сагсыг авах API
func GetCart(c *gin.Context) {
	var cartItems []models.Cart
	userID := c.Param("user_id")

	// JSON форматтай хариу буцаах
	if err := config.DB.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cartItems}) // ✅ JSON форматтай буцаах
}

// ❌ Сагснаас бүтээгдэхүүн устгах
func RemoveFromCart(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.Cart{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove product from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product removed from cart"})
}
