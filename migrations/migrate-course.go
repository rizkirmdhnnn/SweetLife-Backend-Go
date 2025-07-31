package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
)

func main() {
	// Load environment variables and database
	config.LoadEnv()
	config.LoadDatabase()

	// Open the CSV file
	file, err := os.Open("educations.csv")
	if err != nil {
		log.Fatal("Failed to open CSV file:", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)
	reader.LazyQuotes = true // Handle quotes properly

	// Read the header
	header, err := reader.Read()
	if err != nil {
		log.Fatal("Failed to read CSV header:", err)
	}

	fmt.Println("CSV Headers:", header)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to read CSV records:", err)
	}

	fmt.Printf("Found %d records to import\n", len(records))

	// Clear existing mini course data
	if err := config.DB.Exec("DELETE FROM mini_courses").Error; err != nil {
		log.Fatal("Failed to clear existing data:", err)
	}

	fmt.Println("Cleared existing mini course data")

	// Import records
	successCount := 0
	errorCount := 0

	for i, record := range records {
		if len(record) < 5 {
			fmt.Printf("Skipping record %d: insufficient columns\n", i+1)
			errorCount++
			continue
		}

		// Map CSV columns to model fields
		// CSV columns: Title, Image URL, URL, Label, Description
		title := strings.TrimSpace(record[0])
		image := strings.TrimSpace(record[1])
		url := strings.TrimSpace(record[2])
		description := strings.TrimSpace(record[4])

		// Validate required fields
		if title == "" || image == "" || url == "" || description == "" {
			fmt.Printf("Skipping record %d: missing required fields\n", i+1)
			errorCount++
			continue
		}

		// Create mini course record
		miniCourse := models.MiniCourse{
			Title:       title,
			Image:       image,
			URL:         url,
			Description: description,
		}

		// Insert into database
		if err := config.DB.Create(&miniCourse).Error; err != nil {
			fmt.Printf("Failed to insert record %d: %v\n", i+1, err)
			errorCount++
			continue
		}

		successCount++
		if successCount%100 == 0 {
			fmt.Printf("Imported %d records...\n", successCount)
		}
	}

	fmt.Printf("\nMigration completed!\n")
	fmt.Printf("Successfully imported: %d records\n", successCount)
	fmt.Printf("Failed to import: %d records\n", errorCount)
	fmt.Printf("Total processed: %d records\n", len(records))
}
