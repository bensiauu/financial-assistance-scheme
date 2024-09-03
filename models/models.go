package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	Income           int               `gorm:"default:0;not null"`
	MaritalStatus    string            `gorm:"size:50;not null"`
	DisabilityStatus string            `gorm:"size:50;not null"`
	NumberOfChildren int               `gorm:"not null"`
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
	Income           int               `json:"income"`
	MaritalStatus    string            `json:"marital_status"`
	DisabilityStatus string            `json:"disability_status"`
	NumberOfChildren int               `json:"number_of_children"`
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

type Rule struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type Criteria struct {
	Rules []Rule `json:"rules"`
}

func (c *Criteria) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &c)
}

func (c Criteria) Value() (driver.Value, error) {
	return json.Marshal(c)
}

type Scheme struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string          `gorm:"size:255;not null"`   // Name of the scheme
	Criteria  Criteria        `gorm:"type:jsonb;not null"` // Criteria for eligibility (stored as JSONB)
	Benefits  json.RawMessage `gorm:"type:jsonb;not null"` // Benefits provided by the scheme (stored as JSONB)
	CreatedAt time.Time       `gorm:"autoCreateTime"`      // Timestamp of when the scheme was created
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`      // Timestamp of when the scheme was last updated
}

type Application struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ApplicantID uuid.UUID `gorm:"type:uuid;not null"`                 // Foreign key to applicants
	SchemeID    uuid.UUID `gorm:"type:uuid;not null"`                 // Foreign key to schemes
	Status      string    `gorm:"size:50;not null;default:'pending'"` // Status of the application
	CreatedAt   time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
}
