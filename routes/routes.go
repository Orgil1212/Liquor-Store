package routes

import (
	"liquor-store/controllers"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 📌 CORS тохиргоо хийх
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Frontend-ээс зөвшөөрөх
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.GET("/profile/:id", controllers.GetProfile)
		api.POST("/check-user", controllers.CheckUserExists) // 📌 Хэрэглэгч бүртгэлтэй эсэхийг шалгах API
		api.POST("/forgot-password", controllers.ForgotPassword)
		api.GET("/verify-email/:token", controllers.VerifyEmail)   // Токен параметрийг зөвшөөрөх
		api.GET("/products", controllers.GetProducts)              // Бүтээгдэхүүн авах
		api.POST("/products", controllers.CreateProduct)           // Бүтээгдэхүүн нэмэх
		api.PUT("/update-profile", controllers.UpdateProfile)      // 📌 Профайл шинэчлэх API
		api.DELETE("/api/products/:id", controllers.DeleteProduct) //
		api.POST("/api/cart", controllers.AddToCart)               // 🛒 Сагсанд нэмэх
		api.GET("/api/cart/:user_id", controllers.GetCart)         // 📦 Хэрэглэгчийн сагсыг авах
		api.DELETE("/api/cart/:id", controllers.RemoveFromCart)    // ❌ Сагснаас бүтээгдэхүүн хасах

		// Админ API
		admin := api.Group("/admin")
		admin.Use(AdminMiddleware()) // Админ middleware ашиглах
		{
			admin.GET("/users", controllers.GetAllUsers) // Жишээ админ API
		}
	}

	return r
}

// AdminMiddleware - зөвхөн админ хэрэглэгчдийг зөвшөөрнө
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization header-аас токеныг авах
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен байхгүй байна!"})
			c.Abort()
			return
		}

		// "Bearer " үгийг арилгах
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// JWT токеныг шалгах
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Хүчингүй токен!"})
			c.Abort()
			return
		}

		// Токеноос role талбарыг шалгах
		role, ok := (*claims)["role"].(string)
		if !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Админ эрх шаардлагатай!"})
			c.Abort()
			return
		}

		c.Next()
	}
}
