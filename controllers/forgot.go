package controllers

import (
	"liquor-store/config"
	"liquor-store/models"
	"liquor-store/utils"
	"log" // 📌 Санамсаргүй тоо үүсгэх сан
	"net/http"

	"github.com/gin-gonic/gin"
)

func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}

	// Имэйл хүлээн авах
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("❌ Имэйл хүлээн авахад алдаа гарлаа:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email заавал шаардлагатай!"})
		return
	}

	// Хэрэглэгчийн мэдээллийг олох
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		log.Println("❌ Хэрэглэгч олдсонгүй:", err) // Алдааг хэвлэх
		c.JSON(http.StatusNotFound, gin.H{"error": "Хэрэглэгч олдсонгүй!"})
		return
	}

	// Нэг удаагийн баталгаажуулах код үүсгэх
	verificationCode := utils.GenerateVerificationCode()
	log.Println("✅ Баталгаажуулах код:", verificationCode) // Баталгаажуулах код

	// Түр кодыг илгээх
	if err := utils.SendVerificationCodeEmail(user.Email, verificationCode); err != nil {
		log.Println("❌ Имэйл илгээхэд алдаа гарлаа:", err) // Имэйл илгээх алдаа
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Имэйл илгээхэд алдаа гарлаа!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Таны баталгаажуулах код имэйл хаяг руу илгээгдлээ."})
}
