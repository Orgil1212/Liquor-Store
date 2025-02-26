package main

import (
	"fmt"
	"liquor-store/config"
	"liquor-store/routes"
	"liquor-store/utils"
)

func main() {
	code := utils.GenerateVerificationCode()
	fmt.Println("Баталгаажуулах код:", code)
	// Өгөгдлийн сантай холбогдох
	config.ConnectDatabase()
	// Серверийг эхлүүлэх
	r := routes.SetupRouter()
	r.Static("/uploads", "./uploads")
	fmt.Println("🚀 Сервер 8080 порт дээр ажиллаж байна...")
	r.Run(":8080") // Серверийг 8080 порт дээр ажиллуулах
}
