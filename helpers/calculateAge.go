package helper

import (
	"fmt"
	"time"
)

func CalculateAge(birthDate string) (int, error) {
	// Format tanggal lahir
	layout := "2006-01-02"

	// Parsing tanggal lahir ke tipe time.Time
	birthTime, err := time.Parse(layout, birthDate)
	if err != nil {
		return 0, fmt.Errorf("invalid date format: %v", err)
	}

	// Mendapatkan tanggal hari ini
	currentTime := time.Now()

	// Menghitung umur
	age := currentTime.Year() - birthTime.Year()

	// Memastikan bulan dan hari belum terlewati dalam tahun berjalan
	if currentTime.Month() < birthTime.Month() || (currentTime.Month() == birthTime.Month() && currentTime.Day() < birthTime.Day()) {
		age--
	}

	return age, nil
}
