package models

import (
	"time"

	"gorm.io/gorm"
)

type Food struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Volume    string         `gorm:"type:varchar(100)" json:"volume"`
	ImageURL  string         `gorm:"type:text" json:"image_url"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type FoodNutrition struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	FoodID        uint           `gorm:"not null;index" json:"food_id"`
	Food          Food           `json:"food" gorm:"foreignKey:FoodID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Calories      float64        `gorm:"not null" json:"calories"`
	Sugar         float64        `gorm:"not null" json:"sugar"`
	Fat           float64        `gorm:"not null" json:"fat"`
	Carbohydrates float64        `gorm:"not null" json:"carbohydrates"`
	Proteins      float64        `gorm:"not null" json:"proteins"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserFoodHistory struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	User      User           `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FoodID    uint           `gorm:"not null;index" json:"food_id"`
	Food      Food           `json:"food" gorm:"foreignKey:FoodID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Date      time.Time      `gorm:"not null" json:"date"`
	Portion   float64        `gorm:"not null" json:"portion"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ScanFood struct {
	Name        string  `json:"Name"`
	JumlahUnit  int     `json:"Jumlah Unit"`
	Berat       float64 `json:"Berat (g)"`
	Kalori      float64 `json:"Kalori (Kcal)"`
	Protein     float64 `json:"Protein (g)"`
	Gula        float64 `json:"Gula (g)"`
	Karbohidrat float64 `json:"Karbohidrat (g)"`
	Lemak       float64 `json:"Lemak (g)"`
}
