package repositories

import (
	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"gorm.io/gorm"
)

// UserRepository is a contract of user repository
type UserRepository interface {
	Update(user *models.User) error
	GetFoodHistory(userID string) ([]dto.FoodHistoryByDate, error)
	GetDailyNutrition(userID string) (*dto.DailyNutrition, error)
}

// userRepository is a struct to store db connection
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository is a constructor to create user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &userRepository{
		db: db,
	}
}

// Implement Update method of userRepository
func (r *userRepository) Update(user *models.User) error {
	err := r.db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// GetFoodHistory retrieves food history for a user with pagination.
func (r *userRepository) GetFoodHistory(userID string) ([]dto.FoodHistoryByDate, error) {
	var foodHistory []dto.FoodHistoryByDate

	// Get food history data with the given page and page size
	err := r.db.Raw(`
		SELECT 
			user_food_histories.id AS id,
			DATE(user_food_histories.created_at) AS date,
			SUM(user_food_histories.unit) AS total_units,
			foods.name AS food_name,
			food_nutritions.calories * user_food_histories.unit AS calories,
			TO_CHAR(user_food_histories.created_at, 'HH24:MI') AS time
		FROM 
			user_food_histories
		JOIN 
			foods ON foods.id = user_food_histories.food_id
		JOIN 
			food_nutritions ON food_nutritions.food_id = foods.id
		WHERE 
			user_food_histories.user_id = ?
		GROUP BY 
			DATE(user_food_histories.created_at), 
			user_food_histories.id, 
			foods.name, 
			food_nutritions.calories, 
			TO_CHAR(user_food_histories.created_at, 'HH24:MI')
		ORDER BY 
			DATE(user_food_histories.created_at) DESC`, userID).
		Scan(&foodHistory).Error

	if err != nil {
		return nil, err
	}

	return foodHistory, nil
}

// GetDailyNutrition implements UserRepository.
func (r *userRepository) GetDailyNutrition(userID string) (*dto.DailyNutrition, error) {
	var dailyNutrition dto.DailyNutrition
	r.db.Raw(`
	SELECT 
    DATE(user_food_histories.created_at) AS date,
    ROUND(
        SUM(
            CASE 
                WHEN user_food_histories.weight IS NOT NULL AND user_food_histories.weight > 0 
                THEN food_nutritions.calories * (user_food_histories.weight / 100.0)
                ELSE food_nutritions.calories * user_food_histories.unit
            END
        ), 1
    ) AS total_calories,
    ROUND(
        SUM(
            CASE 
                WHEN user_food_histories.weight IS NOT NULL AND user_food_histories.weight > 0 
                THEN food_nutritions.sugar * (user_food_histories.weight / 100.0)
                ELSE food_nutritions.sugar * user_food_histories.unit
            END
        ), 1
    ) AS total_sugar,
    ROUND(
        SUM(
            CASE 
                WHEN user_food_histories.weight IS NOT NULL AND user_food_histories.weight > 0 
                THEN food_nutritions.carbohydrates * (user_food_histories.weight / 100.0)
                ELSE food_nutritions.carbohydrates * user_food_histories.unit
            END
        ), 1
    ) AS total_carbs
FROM 
    user_food_histories
JOIN 
    foods ON foods.id = user_food_histories.food_id
JOIN 
    food_nutritions ON food_nutritions.food_id = foods.id
WHERE 
    user_food_histories.user_id = ?
GROUP BY 
    DATE(user_food_histories.created_at)
ORDER BY 
    DATE(user_food_histories.created_at) DESC;
	`, userID).Scan(&dailyNutrition)

	return &dailyNutrition, nil
}
