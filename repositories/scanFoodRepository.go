package repositories

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
)

type ScanFoodRepository interface {
	ScanFood(image string) (*dto.ScanFoodClientResp, error)
}

type scanFoodRepository struct {
	httpClient *http.Client
}

func NewScanFoodRepository(httpClient *http.Client) ScanFoodRepository {
	return &scanFoodRepository{
		httpClient: httpClient,
	}
}

func (s *scanFoodRepository) ScanFood(image string) (*dto.ScanFoodClientResp, error) {
	data := map[string]interface{}{
		"image": image,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Post("https://mock.apidog.com/m1/739297-716037-3350759e/scan_food", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var scanFoodResponse dto.ScanFoodClientResp
	err = json.NewDecoder(resp.Body).Decode(&scanFoodResponse)
	if err != nil {
		return nil, err
	}

	return &scanFoodResponse, nil
}
