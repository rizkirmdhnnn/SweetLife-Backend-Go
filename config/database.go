package config

import (
	"fmt"
	"log"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDatabase() {
	// Database connection string
	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		ENV.DB_HOST, ENV.DB_USER, ENV.DB_PASSWORD, ENV.DB_NAME, ENV.DB_PORT)

	// Open a database connection
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		log.Fatal("Failed to connect database")
	}

	// Create extension uuid-ossp
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// Migrate the schema
	if err := db.AutoMigrate(
		models.User{},
		models.Password_reset_tokens{},
		models.RefreshToken{},
		models.HealthProfile{},
		models.DiabetesDetails{},
		models.RiskAssessment{},
		models.Food{},
		models.FoodNutrition{},
		models.UserFoodHistory{},
		models.MiniCourse{},
	); err != nil {
		log.Fatal("Failed to migrate table")
	}

	DB = db
	log.Println("Database connected")
}
