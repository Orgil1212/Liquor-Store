package models

import "gorm.io/gorm"

// Сагсны бүтэц
type Cart struct {
	gorm.Model
	UserID    uint `json:"user_id"`    // Хэрэглэгчийн ID
	ProductID uint `json:"product_id"` // Бүтээгдэхүүний ID
	Quantity  int  `json:"quantity"`   // Тоо ширхэг
}
