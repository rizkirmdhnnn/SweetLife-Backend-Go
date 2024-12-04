package dto

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
	Amount       int     `json:"amount"`
	Weight       float64 `json:"weight"`
	Calories     float64 `json:"calories"`
	Protein      float64 `json:"protein"`
	Sugar        float64 `json:"sugar"`
	Carbohydrate float64 `json:"carbohydrate"`
	Fat          float64 `json:"fat"`
}
