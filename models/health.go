package models

import "time"

type ActivityLevel string

const (
	Sedentary ActivityLevel = "never"     // Tidak pernah olahraga
	Light     ActivityLevel = "light"     // Olahraga 1 - 3 hari per minggu
	Moderate  ActivityLevel = "moderate"  // Olahraga 3 - 5 hari per minggu
	Active    ActivityLevel = "active"    // Olahraga 6 - 7 hari per minggu
	Extremely ActivityLevel = "extremely" // setiap hari bisa 2x dalam sehari
)

type DiabeticType string

const (
	Type1       DiabeticType = "type1"
	Type2       DiabeticType = "type2"
	Type3       DiabeticType = "type3"
	Gestational DiabeticType = "gestational"
)

type RiskLevelType string

const (
	LowRisk    RiskLevelType = "low"
	MediumRisk RiskLevelType = "medium"
	HighRisk   RiskLevelType = "high"
)

type SmokingHistory string

const (
	Never   SmokingHistory = "never"
	Current SmokingHistory = "current"
	Former  SmokingHistory = "former"
	Ever    SmokingHistory = "ever"
)

type HealthProfile struct {
	ID              uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          string         `json:"user_id" gorm:"type:uuid;not null;unique;index"`
	User            User           `json:"user" gorm:"foreignKey:UserID"`
	Height          float64        `json:"height" gorm:"not null;type:decimal(5,2)"`
	Weight          float64        `json:"weight" gorm:"not null;type:decimal(5,2)"`
	BMI             float64        `json:"bmi" gorm:"not null;type:decimal(4,2)"`
	IsDiabetic      bool           `json:"is_diabetic" gorm:"not null;default:false"`
	SmokingHistory  SmokingHistory `json:"smoking_history" gorm:"type:varchar(10);not null"`
	HasHeartDisease bool           `json:"has_heart_disease" gorm:"not null;default:false"`
	ActivityLevel   ActivityLevel  `json:"activity_level" gorm:"type:varchar(10);not null"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

type DiabetesDetails struct {
	ID            uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	ProfileID     uint          `json:"profile_id" gorm:"not null;index;unique"`
	Profile       HealthProfile `json:"profile" gorm:"foreignKey:ProfileID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DiabeticType  DiabeticType  `json:"diabetic_type" gorm:"type:varchar(15);not null"`
	InsulinLevel  float64       `json:"insulin_level" gorm:"not null;type:decimal(6,2)"`
	BloodPressure uint          `json:"blood_pressure" gorm:"not null"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

type RiskAssessment struct {
	ID         uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	ProfileID  uint             `json:"profile_id" gorm:"not null;index"`
	Profile    HealthProfile    `json:"profile" gorm:"foreignKey:ProfileID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DiabetesID *uint            `json:"diabetes_id" gorm:"index"`
	Diabetes   *DiabetesDetails `json:"diabetes" gorm:"foreignKey:DiabetesID"`
	RiskLevel  RiskLevelType    `json:"risk_level" gorm:"type:varchar(10);not null"`
	RiskScore  float64          `json:"risk_score" gorm:"not null;type:decimal(5,2)"`
	Note       string           `json:"note" gorm:"type:text"`
	CreatedAt  time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}
