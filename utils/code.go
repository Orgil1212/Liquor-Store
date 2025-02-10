package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// Баталгаажуулах код үүсгэх функц
func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())                // Санамсаргүй тоо үүсгэхийн тулд цаг хугацаа ашиглана
	code := fmt.Sprintf("%06d", rand.Intn(1000000)) // 6 оронтой санамсаргүй код үүсгэнэ
	return code
}
