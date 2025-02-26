package controllers

import (
	"liquor-store/config"
	"liquor-store/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Хэрэглэгчийн профайлыг шинэчлэх API
func UpdateProfile(c *gin.Context) {
	var input struct {
		Address string `json:"address"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
	}

	// Клиентээс ирсэн JSON өгөгдлийг авах
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON өгөгдөл буруу байна!"})
		return
	}

	// Токен-аас хэрэглэгчийн ID авах
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен байхгүй байна!"})
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Хүчингүй токен!"})
		return
	}

	userID, ok := (*claims)["id"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Хэрэглэгчийн ID олдсонгүй!"})
		return
	}

	// Өгөгдлийн сангаас хэрэглэгчийг авах
	var user models.User
	if err := config.DB.First(&user, int(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Хэрэглэгч олдсонгүй!"})
		return
	}

	// Шинэ утгуудыг оноох
	user.Address = input.Address
	user.Email = input.Email
	user.Phone = input.Phone

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Өгөгдөл хадгалахад алдаа гарлаа!"})
		return
	}

	// Хэрэглэгчийн шинэ мэдээллийг буцаах
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully!",
		"user": gin.H{
			"address": user.Address,
			"email":   user.Email,
			"phone":   user.Phone,
		},
	})
}
