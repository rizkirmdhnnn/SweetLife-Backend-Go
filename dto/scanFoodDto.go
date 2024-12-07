package dto

type FoodNutritionResponseClient struct {
	Foods []struct {
		Description   string `json:"description"`
		FoodCategory  string `json:"foodCategory"`
		FoodNutrients []struct {
			NutrientName string  `json:"nutrientName"`
			Value        float64 `json:"value"`
			UnitName     string  `json:"unitName"`
		} `json:"foodNutrients"`
	} `json:"foods"`
}

type FoodNutritionResponse struct {
	Name     string  `json:"name"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
	Sugar    float64 `json:"sugar"`
	Weight   float64 `json:"weight"`
}

type ScanFoodClientResp struct {
	Objects []ScanFood `json:"objects"`
}

type ScanFood struct {
	Name string `json:"name"`
	Unit int    `json:"unit"`
}

type ScanFoodResponse struct {
	IsDetected bool       `json:"is_detected"`
	FoodList   []FoodList `json:"food_list"`
}

type FoodList struct {
	Name         string  `json:"name"`
	Unit         int     `json:"unit"`
	Calories     float64 `json:"calories"`
	Protein      float64 `json:"protein"`
	Sugar        float64 `json:"sugar"`
	Carbohydrate float64 `json:"carbohydrate"`
	Fat          float64 `json:"fat"`
}

type FindFoodClientResp struct {
	Alert         string                  `json:"alert"`
	FoodName      string                  `json:"food_name"`
	NutritionInfo NutritionInfoClientResp `json:"nutrition_info"`
	Weight        float64                 `json:"weight"`
}

type NutritionInfoClientResp struct {
	Calories      float64 `json:"calories"`
	Carbohydrates float64 `json:"carbohydrates"`
	Fat           float64 `json:"fat"`
	Proteins      float64 `json:"proteins"`
	Sugar         float64 `json:"sugar"`
}

type FindFoodRequest struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

type SaveFoodRequest struct {
	Scan        []ScanFood `json:"scan"`
	Additionall []struct {
		Name   string  `json:"name"`
		Weight float64 `json:"weight"`
	} `json:"additionall"`
}
