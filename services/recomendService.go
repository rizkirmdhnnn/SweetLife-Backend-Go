package services

import (
	"errors"
	"fmt"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
	"gorm.io/gorm"
)

type RecomendationService interface {
	GetFoodRecomendations(userid string) ([]*dto.FoodRecomendation, error)
	GetExerciseRecomendations(userid string) (*dto.ExerciseRecommendation, error)
}

type recomendationService struct {
	recomendationRepo repositories.RecomendationRepo
	healthRepo        repositories.HealthProfileRepository
	authRepo          repositories.AuthRepository
}

func NewRecomendationService(recomendationRepo repositories.RecomendationRepo, healthRepo repositories.HealthProfileRepository, authRepo repositories.AuthRepository) RecomendationService {
	if recomendationRepo == nil {
		panic("recomendationRepo cannot be nil")
	}
	return &recomendationService{
		recomendationRepo: recomendationRepo,
		healthRepo:        healthRepo,
		authRepo:          authRepo,
	}
}

// GetRecomendations implements RecomendationService.
func (r *recomendationService) GetFoodRecomendations(userid string) ([]*dto.FoodRecomendation, error) {
	healthProfile, err := r.healthRepo.GetRiskAssessmentByUserID(userid)
	var riskScore float32

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			riskScore = 100
		} else {
			return nil, err
		}
	} else {
		riskScore = float32(healthProfile.RiskScore)
	}

	// 2. Get recommendations
	foodRecomendationClientResp, err := r.recomendationRepo.GetFoodRecomendations(riskScore)
	if err != nil {
		return nil, err
	}
	var foodRecomendations []*dto.FoodRecomendation
	for _, foodList := range foodRecomendationClientResp.FoodRecomendation {
		for _, food := range foodList {
			foodRec := dto.FoodRecomendation{
				Name: food.Name,
				Details: dto.RecomendationDetails{
					Carbohydrate: fmt.Sprintf("%.2f g", food.Carbohydrate),
					Calories:     fmt.Sprintf("%.2f kcal", food.Calories),
					Fat:          fmt.Sprintf("%.2f g", food.Fat),
					Proteins:     fmt.Sprintf("%.2f g", food.Proteins),
				},
				Image: food.Image,
			}
			foodRecomendations = append(foodRecomendations, &foodRec)
		}
	}

	return foodRecomendations, nil
}

// GetExerciseRecomendations implements RecomendationService.
func (r *recomendationService) GetExerciseRecomendations(userid string) (*dto.ExerciseRecommendation, error) {
	// get user health profile
	healthProfile, err := r.healthRepo.GetHealthProfileByUserID(userid)
	if err != nil {
		//TODO: if user not found
		return nil, err
	}

	// get user profile
	userProfile, err := r.authRepo.GetUserById(userid)
	if err != nil {
		return nil, err
	}

	if userProfile.Age == 0 || userProfile.DateOfBirth.IsZero() {
		return nil, fmt.Errorf("please update your profile")
	}

	// create request
	userData := dto.ExerciseRequest{
		Age:      userProfile.Age,
		Height:   healthProfile.Height,
		Bmi:      healthProfile.BMI,
		Gender:   userProfile.Gender,
		Diabetes: healthProfile.IsDiabetic,
	}

	// get recomendation
	exerciseRecomendationClientResp, err := r.recomendationRepo.GetExerciceRecomendations(&userData)
	if err != nil {
		return nil, err
	}

	exerciseData := []dto.ExerciseList{
		{
			Name:  "Squats",
			Desc:  "A strength training exercise involving the thighs, hips, and buttocks muscles. Squats help improve lower body strength and core stability.",
			Image: "https://storage.googleapis.com/sweetlife-go/website/exercise/Squats.jpg",
		},
		{
			Name:  "Deadlifts",
			Desc:  "A weightlifting exercise where you lift a weight from the floor to your hips. Deadlifts strengthen the lower back, thighs, buttocks, and upper back muscles.",
			Image: "https://storage.googleapis.com/sweetlife-go/website/exercise/Deadlifts.jpg",
		},
		{
			Name:  "Bench presses",
			Desc:  "A strength exercise using a barbell or dumbbells, typically performed while lying on a bench. It focuses on strengthening the chest, shoulders, and triceps.",
			Image: "https://storage.googleapis.com/sweetlife-go/website/exercise/Bench%20presses.jpg",
		},
		{
			Name:  "Overhead presses",
			Desc:  "A weightlifting exercise where you push a weight overhead from your shoulders. It helps train the shoulder muscles, triceps, and core stability.",
			Image: "https://storage.googleapis.com/sweetlife-go/website/exercise/Overhead%20presses.jpg",
		},
		{
			Name:  "Yoga",
			Desc:  "A physical and mental practice combining body postures, breathing techniques, and meditation. Yoga improves flexibility, balance, strength, and reduces stress.",
			Image: "https://storage.googleapis.com/sweetlife-go/website/exercise/Yoga.jpg",
		},
		{
			Name:  "Brisk walking",
			Desc:  "A fast-paced walk aimed at increasing heart rate. It's beneficial for heart health, calorie burning, and general fitness.",
			Image: "https://storage.googleapis.com/sweetlife-go/website/exercise/Brisk%20walking.jpg",
		},
		{
			Name:  "Cycling",
			Desc:  "An activity involving pedaling a bicycle that works the leg muscles, strengthens cardiovascular health, and burns calories. It can be done outdoors or on a stationary bike.",
			Image: "https://storage.googleapis.com/sweetlife-go/website/exercise/Cycling.jpg",
		},
		{
			Name:  "Swimming",
			Desc:  "A water sport that involves almost all the body's muscles. It helps improve endurance, breathing techniques, and protects joints due to low impact.",
			Image: "https://storage.googleapis.com/sweetlife-go/website/exercise/Swimming.jpg",
		},
		{
			Name:  "Running",
			Desc:  "A cardiovascular exercise that helps improve heart health, endurance, and burns a significant amount of calories.",
			Image: "https://storage.googleapis.com/sweetlife-go/website/exercise/Running.jpg",
		},
		{
			Name:  "Dancing",
			Desc:  "A physical activity involving rhythmic body movements to music. It's great for fitness, coordination, and mood improvement.",
			Image: "https://storage.googleapis.com/sweetlife-go/website/exercise/Dancing.jpg",
		},
		{
			Name:  "Walking",
			Desc:  "A light activity that can be done by anyone. It helps improve blood circulation, reduce stress, and maintain heart health.",
			Image: "https://storage.googleapis.com/sweetlife-go/website/exercise/Walking.jpg",
		},
	}

	// Membuat map dari exerciseData untuk kemudahan pencarian berdasarkan nama
	exerciseDataMap := make(map[string]dto.ExerciseList) // Perbaiki tipe data map
	for _, exercise := range exerciseData {
		exerciseDataMap[exercise.Name] = exercise
	}

	// List hasil gabungan
	var exerciseList []*dto.ExerciseList

	// Loop untuk menggabungkan data dari ExerciseCategories dengan ExerciseData
	for _, nameExercise := range exerciseRecomendationClientResp.ExerciseCategories {
		// Membuat struct baru untuk ExerciseList
		exeList := dto.ExerciseList{
			Name: nameExercise,
		}

		// Jika nama olahraga ditemukan di map, ambil deskripsi dan gambar
		if exercise, found := exerciseDataMap[nameExercise]; found {
			exeList.Desc = exercise.Desc
			exeList.Image = exercise.Image
		} else {
			// Jika nama olahraga tidak ada di exerciseData, tentukan default atau kosongkan
			exeList.Desc = ""
			exeList.Image = ""
		}

		// Menambahkan hasil gabungan ke dalam list
		exerciseList = append(exerciseList, &exeList)
	}

	exerciseRecomendations := dto.ExerciseRecommendation{
		CaloriesBurned:   exerciseRecomendationClientResp.CaloriesBurned,
		ExerciseDuration: exerciseRecomendationClientResp.ExerciseDuration,
		ExerciseList:     exerciseList,
	}

	return &exerciseRecomendations, nil
}
