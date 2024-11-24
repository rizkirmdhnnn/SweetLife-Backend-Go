package models

import (
	"time"
)

// thanks copilot for generating documentation

// User model
// User represents a user in the system.
// It contains personal information such as name, email, and date of birth,
// as well as metadata like creation and update timestamps.
//
// Fields:
// - ID: Unique identifier for the user, generated as a UUID.
// - Name: Full name of the user.
// - Email: Email address of the user, must be unique and not null.
// - Password: Hashed password of the user.
// - DateOfBirth: Date of birth of the user.
// - Gender: Gender of the user.
// - ImageUrl: URL to the user's profile image, can be null.
// - Verified_at: Timestamp when the user's email was verified, can be null.
// - Created_at: Timestamp when the user was created.
// - Updated_at: Timestamp when the user was last updated.
type User struct {
	ID          string     ` json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string     `json:"name" gorm:"type:varchar(255)"`
	Email       string     `json:"email" gorm:"type:varchar(100);uniqueIndex; not null"`
	Password    string     `json:"password" gorm:"type:varchar(255)"`
	DateOfBirth time.Time  `gorm:"type:date" json:"date_of_birth"`
	Gender      string     `json:"gender" gorm:"type:varchar(10)"`
	ImageUrl    string     `json:"image_url" gorm:"type:text;default:null"`
	Verified_at *time.Time `json:"verified_at" gorm:"default:null"`
	Created_at  time.Time  `json:"created_at"`
	Updated_at  time.Time  `json:"updated_at"`
}

// Password_reset_tokens model
// Password_reset_tokens represents the structure for password reset tokens.
// It includes fields for the token ID, associated email, the token itself, and its expiration time.
type Password_reset_tokens struct {
	ID         string    ` json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Email      string    `json:"email" gorm:"type:varchar(100);uniqueIndex; not null"`
	Token      string    `json:"token" gorm:"type:varchar(255)"`
	Expires_at time.Time `json:"expires_at"`
}

// RefreshToken model
// RefreshToken represents a refresh token model in the system.
// It includes information about the token itself, the user it belongs to,
// its expiration time, creation and update timestamps, and whether it has been revoked.
// Fields:
// - ID: Primary key, auto-incremented.
// - UserID: UUID of the user associated with the token, indexed.
// - User: The user associated with the token.
// - RefreshToken: The actual refresh token string, indexed.
// - ExpiresAt: The expiration time of the token, indexed.
// - CreatedAt: The timestamp when the token was created, auto-generated.
// - Revoked: Indicates whether the token has been revoked, default is false, indexed.
// - DeviceInfo: Information about the device from which the token was generated.
// - UpdatedAt: The timestamp when the token was last updated, auto-generated.
type RefreshToken struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       string    `gorm:"type:uuid;not null;index" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
	RefreshToken string    `gorm:"type:text;not null;index" json:"refresh_token"`
	ExpiresAt    time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	Revoked      bool      `gorm:"default:false;index" json:"revoked"`
	DeviceInfo   string    `gorm:"type:text" json:"device_info"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
