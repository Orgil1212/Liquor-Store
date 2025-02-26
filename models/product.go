package models

// Product загвар нь бүтээгдэхүүний мэдээллийг хадгална
type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey"` // Бүтээгдэхүүний ID
	Name        string  `json:"name"`                 // Бүтээгдэхүүний нэр
	Price       float64 `json:"price"`                // Бүтээгдэхүүний үнэ
	Image       string  `json:"image"`                // Бүтээгдэхүүний зураг
	Description string  `json:"description"`          // Бүтээгдэхүүний тодорхойлолт
}
