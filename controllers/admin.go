package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization header-аас токеныг авах
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Токен байхгүй байна!"})
			c.Abort()
			return
		}

		// "Bearer " текстийг устгах
		tokenString = tokenString[7:]

		// JWT токеныг шалгах
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Хүчингүй токен!"})
			c.Abort()
			return
		}

		// Role талбарыг шалгах
		role, ok := (*claims)["role"].(string)
		if !ok || role != "admin" {
			c.JSON(403, gin.H{"error": "Энэ үйлдлийг хийх эрхгүй байна!"})
			c.Abort()
			return
		}

		c.Next() // Middleware-г дамжих
	}
}
