package models

import (
	"time"

	"github.com/google/uuid"
)

// Administrator represents a user managing the system.
type Administrator struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name         string    `gorm:"size:255;not null"`
	Email        string    `gorm:"size:255;unique;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// Applicant represents an individual applying for financial assistance.
type Applicant struct {
	ID               uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name             string            `gorm:"size:255;not null"`
	EmploymentStatus string            `gorm:"size:50;not null"`
	Sex              string            `gorm:"size:10;not null"`
	DateOfBirth      time.Time         `gorm:"not null"`
	LastEmployed     *time.Time        `gorm:"type:date"` // Nullable date field
	CreatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP"`
	Household        []HouseholdMember `gorm:"foreignkey:ApplicantID"` // One-to-many relationship
}

// HouseholdMember represents a member of the applicant's household.
type HouseholdMember struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ApplicantID      uuid.UUID `gorm:"type:uuid;not null"` // Foreign key to Applicant
	Name             string    `gorm:"size:255;not null"`
	Relation         string    `gorm:"size:50;not null"` // Relationship to the applicant
	DateOfBirth      time.Time `gorm:"not null"`
	EmploymentStatus string    `gorm:"size:50;not null"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
