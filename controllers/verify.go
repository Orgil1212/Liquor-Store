package controllers

import (
	"liquor-store/config"
	"liquor-store/models" // 📌 Санамсаргүй тоо үүсгэх сан
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Токен байхгүй байна!"})
		return
	}

	var user models.User
	if err := config.DB.Where("verification_token = ?", token).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Хэрэглэгч олдсонгүй!"})
		return
	}

	// Хэрэглэгчийн имэйлийг баталгаажуулах
	user.Verified = true
	user.VerificationToken = ""
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Имэйл амжилттай баталгаажлаа!"})
}

func VerifyCode(c *gin.Context) {
	var input struct {
		Email            string `json:"email"`
		VerificationCode string `json:"verification_code"`
	}

	// Имэйл болон кодыг авах
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Бүх талбарууд шаардлагатай!"})
		return
	}

	// Хэрэглэгчийн мэдээллийг олох
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Хэрэглэгч олдсонгүй!"})
		return
	}

	// Үүсгэсэн кодтой харьцуулах
	if user.VerificationCode != input.VerificationCode {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Төрөл код буруу байна!"})
		return
	}

	// Баталгаажуулалт
	user.Verified = true
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Имэйл амжилттай баталгаажсан!"})
}
