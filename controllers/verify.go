package controllers

import (
	"fmt"
	"liquor-store/config"
	"liquor-store/models" // 📌 Санамсаргүй тоо үүсгэх сан
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func VerifyEmail(c *gin.Context) {
	tokenString := c.Param("token") // Параметрийг авч байна

	// Токен шалгах
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err != nil || !token.Valid {
		c.JSON(400, gin.H{"error": "Хүчингүй баталгаажуулах линк!"})
		return
	}

	email, ok := (*claims)["email"].(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Имэйл мэдээлэл олдсонгүй!"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "Хэрэглэгч олдсонгүй!"})
		return
	}

	user.Verified = true
	config.DB.Save(&user)
	if err != nil {
		fmt.Println("Token Error: ", err)
	}
	if !token.Valid {
		fmt.Println("Invalid Token!")
	}

	c.JSON(200, gin.H{"message": "Баталгаажуулалт амжилттай!"})
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
