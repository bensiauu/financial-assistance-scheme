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

type AdministratorResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Applicant represents an individual applying for financial assistance.
type Applicant struct {
	ID               uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name             string            `gorm:"size:255;not null"`
	EmploymentStatus string            `gorm:"column:employment_status;size:50;not null"`
	Sex              string            `gorm:"size:10;not null"`
	DateOfBirth      time.Time         `gorm:"not null"`
	LastEmployed     *time.Time        `gorm:"type:date"` // Nullable date field
	CreatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP"`
	Household        []HouseholdMember `gorm:"foreignkey:ApplicantID"` // One-to-many relationship
}

type ApplicantResponse struct {
	ID               uuid.UUID         `json:"id"`
	Name             string            `json:"name"`
	EmploymentStatus string            `json:"employment_status"`
	Sex              string            `json:"sex"`
	DateOfBirth      time.Time         `json:"date_of_birth"`
	LastEmployed     *time.Time        `json:"last_employed"`
	Household        []HouseholdMember `json:"household"`
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

type Scheme struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string    `gorm:"size:255;not null"`   // Name of the scheme
	Criteria  string    `gorm:"type:jsonb;not null"` // Criteria for eligibility (stored as JSONB)
	Benefits  string    `gorm:"type:jsonb;not null"` // Benefits provided by the scheme (stored as JSONB)
	CreatedAt time.Time `gorm:"autoCreateTime"`      // Timestamp of when the scheme was created
	UpdatedAt time.Time `gorm:"autoUpdateTime"`      // Timestamp of when the scheme was last updated
}

// Application represents an application for a financial assistance scheme.
type Application struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ApplicantID uuid.UUID `gorm:"type:uuid;not null"`                                                   // Foreign key to Applicant
	Applicant   Applicant `gorm:"foreignkey:ApplicantID;constraint:onUpdate:CASCADE,onDelete:CASCADE;"` // Relationship with Applicant
	SchemeID    uuid.UUID `gorm:"type:uuid;not null"`                                                   // Foreign key to Scheme
	Scheme      Scheme    `gorm:"foreignkey:SchemeID;constraint:onUpdate:CASCADE,onDelete:CASCADE;"`    // Relationship with Scheme
	Status      string    `gorm:"size:50;not null"`                                                     // Status of the application (e.g., pending, approved, rejected)
	CreatedAt   time.Time `gorm:"autoCreateTime"`                                                       // Timestamp of when the application was created
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`                                                       // Timestamp of when the application was last updated
}
