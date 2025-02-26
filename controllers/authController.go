package controllers

import (
	"encoding/hex"
	"fmt"
	"liquor-store/config"
	"liquor-store/models"
	"math/rand" // 📌 Санамсаргүй тоо үүсгэх сан
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

const (
	Argon2Time      uint32 = 3         // Итерацийн тоо
	Argon2Memory    uint32 = 64 * 1024 // 64MB санах ой
	Argon2Threads   uint8  = 4         // Зэрэгцээ утас
	Argon2KeyLength uint32 = 32        // Хэшийн урт (байт)
)

// Хэрэглэгчийн мэдээллийг авах JWT хамгаалагдсан API
func GetProfile(c *gin.Context) {
	// HTTP header-аас токеныг авах
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Токен байхгүй байна!"})
		return
	}

	// "Bearer " гэсэн эхний 7 тэмдэгтийг хасах
	tokenString = tokenString[7:]

	// JWT токеныг шалгах
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Хүчингүй токен!"})
		return
	}

	// Токенээс хэрэглэгчийн ID авах
	userID, ok := (*claims)["id"].(float64)
	if !ok {
		c.JSON(400, gin.H{"error": "Хэрэглэгчийн ID олдсонгүй!"})
		return
	}

	// Өгөгдлийн сангаас хэрэглэгчийн мэдээллийг авах
	var user models.User
	if err := config.DB.First(&user, int(userID)).Error; err != nil {
		c.JSON(404, gin.H{"error": "Хэрэглэгч олдсонгүй!"})
		return
	}

	// Хэрэглэгчийн мэдээллийг буцаах
	c.JSON(200, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}

func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// JSON өгөгдлийг авах
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "JSON өгөгдөл буруу байна!"})
		return
	}

	// Нууц үг хамгийн багадаа 8 тэмдэгттэй байх ёстой
	if len(input.Password) < 8 {
		c.JSON(400, gin.H{"error": "Нууц үг хамгийн багадаа 8 тэмдэгттэй байх ёстой!"})
		return
	}

	// Имэйл давхцахгүй эсэхийг шалгах
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(400, gin.H{"error": "Энэ имэйл аль хэдийн бүртгэгдсэн байна!"})
		return
	}

	// Хэрэглэгчийн мэдээллийг бэлтгэх
	verificationToken := fmt.Sprintf("%x", time.Now().UnixNano()) // Токен үүсгэх

	user := models.User{
		Name:              input.Name,
		Email:             input.Email,
		Password:          input.Password,
		Verified:          false,
		VerificationToken: verificationToken, // Токеныг хадгалах
	}

	// Нууц үгийг шифрлэх
	hashedPassword, salt, err := hashPassword(input.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Нууц үгийг шифрлэхэд алдаа гарлаа!"})
		return
	}
	user.Password = hashedPassword
	user.Salt = salt

	// Өгөгдлийн санд хадгалах
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Өгөгдлийн санд хадгалах үед алдаа гарлаа!"})
		return
	}

	c.JSON(200, gin.H{
		"message":            "Бүртгэл амжилттай! Имэйлээ шалгана уу.",
		"verification_token": user.VerificationToken,
	})
}

// JWT үүсгэх функц

func hashPassword(password string) (string, string, error) {
	// Санамсаргүй salt үүсгэх
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", "", err
	}

	// Argon2 хэш үүсгэх
	hash := argon2.IDKey([]byte(password), salt, Argon2Time, Argon2Memory, Argon2Threads, Argon2KeyLength)

	// Debug мэдээлэл
	fmt.Printf("✅ HashPassword:\n - Password: %s\n - Salt: %x\n - Hash: %x\n", password, salt, hash)

	return fmt.Sprintf("%x", hash), fmt.Sprintf("%x", salt), nil
}

func comparePassword(hashedPassword, salt, password string) bool {
	// Salt хөрвүүлэх
	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		fmt.Println("❌ Salt хөрвүүлэхэд алдаа гарлаа:", err)
		return false
	}

	// Argon2 хэш үүсгэх
	newHash := argon2.IDKey([]byte(password), saltBytes, Argon2Time, Argon2Memory, Argon2Threads, Argon2KeyLength)

	// Debug мэдээлэл
	fmt.Printf("🔍 ComparePassword:\n - Password: %s\n - Salt: %x\n - HashedPassword: %s\n - NewHash: %x\n", password, saltBytes, hashedPassword, newHash)

	return hashedPassword == fmt.Sprintf("%x", newHash)
}

// Нэвтрэх функц
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("❌ JSON Parse Error:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		fmt.Println("❌ User not found")
		c.JSON(401, gin.H{"error": "Хэрэглэгч олдсонгүй!"})
		return
	}

	// Нууц үгийг шалгах
	if !comparePassword(user.Password, user.Salt, input.Password) {
		fmt.Println("❌ Incorrect password")
		c.JSON(401, gin.H{"error": "Нууц үг буруу байна!"})
		return
	}

	// JWT үүсгэх
	token, err := GenerateJWT(user)
	if err != nil {
		fmt.Println("❌ JWT generation failed:", err)
		c.JSON(500, gin.H{"error": "JWT үүсгэхэд алдаа гарлаа!"})
		return
	}

	// ✅ Debug Log - JSON Response харах
	responseData := gin.H{"token": token, "message": "Login successful"}
	fmt.Println("✅ Sending JSON Response:", responseData)

	// ✅ JSON өгөгдлийг зөв буцаах
	c.JSON(200, responseData)
}

func GetAllUsers(c *gin.Context) {
	var users []models.User

	// Бүх хэрэглэгчийн өгөгдлийг авах
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": "Хэрэглэгчдийн жагсаалтыг авахад алдаа гарлаа!"})
		return
	}

	// Хэрэглэгчдийн жагсаалтыг буцаах
	c.JSON(200, gin.H{"users": users})
}

func GenerateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role, // Role талбар ашиглаж байгаа хэсэг
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte("your-secret-key"))
}
func CheckUserExists(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email заавал шаардлагатай!"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"exists": false, "message": "Хэрэглэгч бүртгэлгүй байна."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": true, "message": "Хэрэглэгч бүртгэлтэй байна."})
}

func GenerateTemporaryPassword() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
