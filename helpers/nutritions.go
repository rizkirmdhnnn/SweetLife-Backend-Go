package helper

import "github.com/rizkirmdhnnn/sweetlife-backend-go/models"

// CalculateNutrients menghitung nilai nutrisi berdasarkan berat baru.
func CalculateNutrients(newWeight float64, nutrisi *models.FoodNutrition) models.FoodNutrition {
	ratio := newWeight / 100

	return models.FoodNutrition{
		Calories:      nutrisi.Calories * ratio,
		Sugar:         nutrisi.Sugar * ratio,
		Fat:           nutrisi.Fat * ratio,
		Carbohydrates: nutrisi.Carbohydrates * ratio,
		Proteins:      nutrisi.Proteins * ratio,
		Weight:        newWeight,
	}
}
