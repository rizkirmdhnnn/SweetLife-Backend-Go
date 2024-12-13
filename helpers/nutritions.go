package helper

import (
	"errors"
	"math"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
)

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

func CalculateDailyCalories(req dto.DailyCaloriesRequest) (float64, error) {
	if req.Weight <= 0 || req.Height <= 0 || req.Age <= 0 {
		return 0, errors.New("invalid input: weight, height, and age must be positive numbers")
	}
	if req.Gender != "Male" || req.Gender != "Female" {
		return 0, errors.New("invalid gender: must be 'male' or 'female'")
	}

	// BMR Laki-laki = 66 + (13,7 x BB) + (5 x TB) – (6,78 x U).
	// BMR Perempuan = 655 + (9,6 x BB) + (1,8 x TB) – (4,7 x U).
	// Source : https://eprints.ums.ac.id/78765/3/mufid_Naskah%20Publikasi-143.pdf
	var bmr float64
	if req.Gender == "male" {
		bmr = 66 + (13.7 * req.Weight) + (5 * req.Height) - (6.78 * float64(req.Age))
	} else { // female
		bmr = 655 + (9.6 * req.Weight) + (1.8 * req.Height) - (4.7 * float64(req.Age))
	}

	const (
		Sedentary  = 1.2
		Light      = 1.375
		Moderate   = 1.55
		VeryActive = 1.725
		Extremely  = 1.9
	)

	// Activity level
	switch req.ActivityLevel {
	case models.Sedentary:
		bmr *= Sedentary
	case models.Light:
		bmr *= Light
	case models.Moderate:
		bmr *= Moderate
	case models.Active:
		bmr *= VeryActive
	case models.Extremely:
		bmr *= Extremely
	default:
		return 0, errors.New("invalid activity level")
	}
	multiplier := math.Pow(10, 2)
	return math.Round(bmr*multiplier) / multiplier, nil
}

// https://fkm.unair.ac.id/pentingnya-menjaga-asupan-gula-harian-tubuh/
// https://news.unair.ac.id/2021/05/31/pentingnya-menjaga-asupan-gula-harian-tubuh/?fbclid=IwAR3zVI4olps1v9m4xgG-ydzjQtaZifucTWYEQsqbpjFaXjNwHpXUyaPTHXA
// katanya kalo diabetes maksimal 5% dari total kalori harian, kalo non diabetes 10%
func CalculateDailySugar(dailyCalories float64, isDiabet bool) float64 {
	var result float64
	if dailyCalories <= 0 {
		return 0
	}

	if isDiabet {
		result = dailyCalories * 5 / 100 // 5% untuk diabetes
	} else {
		result = dailyCalories * 10 / 100 // 10% untuk non-diabetes
	}
	multiplier := math.Pow(10, 2)
	return math.Round(result*multiplier) / multiplier
}

// https://www.medicalnewstoday.com/articles/317662#carbs-and-diabetes
func CalculateDialyCarbs(dailyCalories float64) float64 {
	var result float64
	if dailyCalories <= 0 {
		return 0
	}
	result = (dailyCalories * 0.45) / 4
	multiplier := math.Pow(10, 2)
	return math.Round(result*multiplier) / multiplier
}

func DetermineSatisfication(current, target float64) dto.Satisfication {
	percent := (current / target) * 100

	switch {
	case percent < 80:
		return dto.UNDER
	case percent >= 80 && percent <= 120:
		return dto.PASS
	default:
		return dto.OVER
	}
}

func DetermineOverallSatisfication(calories, carbs, sugar dto.DailyProgessPerItem) dto.Satisfication {
	// Buat slice status untuk mempermudah perhitungan
	statuses := []dto.Satisfication{
		calories.Satisfication,
		carbs.Satisfication,
		sugar.Satisfication,
	}

	// Hitung jumlah masing-masing status
	statusCount := map[dto.Satisfication]int{
		dto.UNDER: 0,
		dto.PASS:  0,
		dto.OVER:  0,
	}

	for _, status := range statuses {
		statusCount[status]++
	}

	// Logic penentuan status keseluruhan
	switch {
	case statusCount[dto.OVER] > 0:
		// Jika ada satu parameter OVER, maka keseluruhan dianggap OVER
		return dto.OVER
	case statusCount[dto.UNDER] > 1:
		// Jika lebih dari satu parameter UNDER, maka keseluruhan UNDER
		return dto.UNDER
	default:
		// Dalam kondisi lain, dianggap PASS
		return dto.PASS
	}
}

func DetermineOverallMessage(calories, carbs, sugar dto.Satisfication) string {
	switch {
	case calories == dto.OVER || carbs == dto.OVER || sugar == dto.OVER:
		return "Attention! Some nutrients exceed the limit"
	case calories == dto.UNDER && carbs == dto.UNDER && sugar == dto.UNDER:
		return "You need to increase your nutrient intake"
	default:
		return "You are in good nutritional condition"
	}
}
