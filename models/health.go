package models

import "time"

type ActivityLevel string

const (
	LowActivity      ActivityLevel = "low"
	ModerateActivity ActivityLevel = "moderate"
	HighActivity     ActivityLevel = "high"
)

type DiabeticType string

const (
	Type1       DiabeticType = "type1"
	Type2       DiabeticType = "type2"
	Gestational DiabeticType = "gestational"
)

type RiskLevelType string

const (
	LowRisk    RiskLevelType = "low"
	MediumRisk RiskLevelType = "medium"
	HighRisk   RiskLevelType = "high"
)

type HealthProfile struct {
	ID              uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          string        `json:"user_id" gorm:"type:uuid;not null;index"`
	User            User          `json:"user" gorm:"foreignKey:UserID"`
	Height          float64       `json:"height" gorm:"not null;type:decimal(5,2)"`
	Weight          float64       `json:"weight" gorm:"not null;type:decimal(5,2)"`
	BMI             float64       `json:"bmi" gorm:"not null;type:decimal(4,2)"`
	IsDiabetic      bool          `json:"is_diabetic" gorm:"not null;default:false"`
	IsSmoker        bool          `json:"is_smoker" gorm:"not null;default:false"`
	HasHeartDisease bool          `json:"has_heart_disease" gorm:"not null;default:false"`
	ActivityLevel   ActivityLevel `json:"activity_level" gorm:"type:varchar(10);not null"`
	CreatedAt       time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

type DiabetesDetails struct {
	ID            uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	ProfileID     uint          `json:"profile_id" gorm:"not null;index"`
	Profile       HealthProfile `json:"profile" gorm:"foreignKey:ProfileID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DiabeticType  DiabeticType  `json:"diabetic_type" gorm:"type:varchar(15);not null"`
	InsulinLevel  float64       `json:"insulin_level" gorm:"not null;type:decimal(6,2)"`
	BloodPressure uint          `json:"blood_pressure" gorm:"not null"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

type RiskAssessment struct {
	ID         uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	ProfileID  uint            `json:"profile_id" gorm:"not null;index"`
	Profile    HealthProfile   `json:"profile" gorm:"foreignKey:ProfileID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DiabetesID uint            `json:"diabetes_id" gorm:"index"`
	Diabetes   DiabetesDetails `json:"diabetes" gorm:"foreignKey:DiabetesID"`
	RiskLevel  RiskLevelType   `json:"risk_level" gorm:"type:varchar(10);not null"`
	RiskScore  float64         `json:"risk_score" gorm:"not null;type:decimal(5,2)"`
	CreatedAt  time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}
