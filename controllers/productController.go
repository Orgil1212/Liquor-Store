package controllers

import (
	"fmt"
	"liquor-store/config"
	"liquor-store/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Бүтээгдэхүүн авах API
func GetProducts(c *gin.Context) {
	var products []models.Product

	// 📌 Өгөгдлийн сангаас бүх бүтээгдэхүүнийг авах
	result := config.DB.Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 📌 Хэрэв бүтээгдэхүүн олдоогүй бол 404 буцаах
	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No products found"})
		return
	}

	// ✅ JSON форматтай хариу буцаах
	c.JSON(http.StatusOK, products)
}

func CreateProduct(c *gin.Context) {
	// 1. **Log content-type**
	fmt.Println("Content-Type:", c.Request.Header.Get("Content-Type"))

	// 2. **Multipart form задлах**
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB хүртэл файлын хэмжээ зөвшөөрөх
	if err != nil {
		fmt.Println("Error parsing form:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// 3. **Зөв аргаар өгөгдөл унших (`c.Request.FormValue`)**
	name := c.Request.FormValue("name")
	priceStr := c.Request.FormValue("price")
	description := c.Request.FormValue("description")

	// 4. **Шалгах**
	fmt.Println("Received name:", name)
	fmt.Println("Received price:", priceStr)
	fmt.Println("Received description:", description)

	// 5. **Price хөрвүүлэх**
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
		return
	}

	// 6. **Файл авах**
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println("Image file not found:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	// 7. **Файлыг сервер дээр хадгалах**
	filePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// 8. **Шинэ бүтээгдэхүүн үүсгэх**
	product := models.Product{
		Name:        name,
		Price:       price,
		Description: description,
		Image:       filePath,
	}

	// 9. **Шаардлагатай талбаруудыг шалгах**
	if product.Name == "" || product.Price == 0 || product.Description == "" || product.Image == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// 10. **Өгөгдлийн санд бүтээгдэхүүн нэмэх**
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
		return
	}

	// 11. **Амжилттай хариу буцаах**
	c.JSON(http.StatusOK, gin.H{"message": "Product added successfully", "product": product})
}

func DeleteProduct(c *gin.Context) {
	// URL-аас `id` авах
	id := c.Param("id")

	// Өгөгдлийн санд тухайн бүтээгдэхүүн байгаа эсэхийг шалгах
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Бүтээгдэхүүнийг устгах
	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
