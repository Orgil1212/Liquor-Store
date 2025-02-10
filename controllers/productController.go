package controllers

import (
	"liquor-store/config"
	"liquor-store/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Бүтээгдэхүүн авах API
func GetProducts(c *gin.Context) {
	var products []models.Product

	result := config.DB.Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Бүтээгдэхүүн нэмэх API
func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := config.DB.Create(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Бүтээгдэхүүн нэмэгдлээ!", "product": product})
}
