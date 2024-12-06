package models

import (
	"time"
)

type Food struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type FoodNutrition struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	FoodID        uint      `gorm:"not null;index" json:"food_id"`
	Food          Food      `json:"food" gorm:"foreignKey:FoodID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Calories      float64   `gorm:"not null" json:"calories"`
	Sugar         float64   `gorm:"not null" json:"sugar"`
	Fat           float64   `gorm:"not null" json:"fat"`
	Carbohydrates float64   `gorm:"not null" json:"carbohydrates"`
	Proteins      float64   `gorm:"not null" json:"proteins"`
	Weight        float64   `gorm:"not null" json:"weight"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type FoodWithNutritions struct {
	Food      Food          `json:"food"`
	Nutrition FoodNutrition `json:"nutrition"`
}

type UserFoodHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FoodID    uint      `gorm:"not null;index" json:"food_id"`
	Food      Food      `json:"food" gorm:"foreignKey:FoodID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Unit      int       `gorm:"not null" json:"unit"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ScanFood struct {
	Name          string  `json:"Name"`
	NameIndo      *string `json:"Name_Indo,omitempty"`
	Weight        float64 `json:"Weight"`
	Calories      float64 `json:"Calories"`
	Protein       float64 `json:"Protein"`
	Sugar         float64 `json:"Sugar"`
	Carbohydrates float64 `json:"Carbohydrates"`
	Fat           float64 `json:"Fat"`
}
