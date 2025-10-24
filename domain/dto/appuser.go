package dto

import (
	"time"

	"remedymate-backend/domain/entities"
	"remedymate-backend/util/encrypt"
)

// Consent payloads
type ConsentDTO struct {
	Baseline   []entities.ConsentItem `json:"baseline_consents" binding:"required,dive"`
	Optional   []entities.ConsentItem `json:"optional_consents" binding:"dive"`
	AuditTrail []entities.ConsentItem `json:"audit_trail,omitempty"`
}

// Demographics payloads with ISO/E.164 expectations
type DemographicsDTO struct {
	FullName         encrypt.SecureString       `json:"full_name" binding:"required"`
	DateOfBirth      string                     `json:"date_of_birth" binding:"required"` // YYYY-MM-DD (validate age >= 13)
	Gender           string                     `json:"gender" binding:"required"`
	Height           *entities.Measure          `json:"height,omitempty"`
	Weight           *entities.Measure          `json:"weight,omitempty"`
	BloodType        string                     `json:"blood_type,omitempty"`
	Contact          *entities.ContactInfo      `json:"contact,omitempty"`
	Location         *entities.Location         `json:"location,omitempty"`
	PrimaryLanguage  string                     `json:"primary_language,omitempty"` // ISO 639-1
	EmergencyContact *entities.EmergencyContact `json:"emergency_contact,omitempty"`
}

type MedicalHistoryDTO struct {
	ExistingConditions []entities.ConditionEntry `json:"existing_conditions" binding:"dive"`
	PastSurgeries      []entities.SurgeryEntry   `json:"past_surgeries" binding:"dive"`
	CurrentMedications []entities.Medication     `json:"current_medications" binding:"dive"`
	Allergies          []entities.Allergy        `json:"allergies" binding:"dive"`
	FamilyHistory      []entities.FamilyHistory  `json:"family_history" binding:"dive"`
}

type LifestyleDTO struct {
	DietaryHabits    *entities.DietaryHabits    `json:"dietary_habits,omitempty"`
	PhysicalActivity *entities.PhysicalActivity `json:"physical_activity,omitempty"`
	SubstanceUse     *entities.SubstanceUse     `json:"substance_use,omitempty"`
	Sleep            *entities.Sleep            `json:"sleep,omitempty"`
	Stress           *entities.Stress           `json:"stress,omitempty"`
	Occupation       string                     `json:"occupation,omitempty"`
}

// Response: Onboarding Progress
type OnboardingStatusDTO struct {
	ConsentCompleted        bool      `json:"consent_completed"`
	DemographicsCompleted   bool      `json:"demographics_completed"`
	MedicalHistoryCompleted bool      `json:"medical_history_completed"`
	LifestyleCompleted      bool      `json:"lifestyle_completed"`
	LastUpdated             time.Time `json:"last_updated"`
}
