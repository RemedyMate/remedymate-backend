package entities

import (
	"time"

	"remedymate-backend/util/encrypt"
)

// AppUser stores onboarding/PHI data per user (separate collection from users)
type AppUser struct {
	ID        string           `bson:"_id,omitempty" json:"id"`
	UserID    string           `bson:"userId" json:"user_id"` // unique index
	Status    OnboardingStatus `bson:"status" json:"status"`
	Consent   *ConsentEnvelope `bson:"consent,omitempty" json:"consent,omitempty"`
	Demog     *Demographics    `bson:"demographics,omitempty" json:"demographics,omitempty"`
	History   *MedicalHistory  `bson:"medical_history,omitempty" json:"medical_history,omitempty"`
	Lifestyle *Lifestyle       `bson:"lifestyle,omitempty" json:"lifestyle,omitempty"`
	CreatedAt time.Time        `bson:"createdAt" json:"created_at"`

	UpdatedAt time.Time `bson:"updatedAt" json:"updated_at"`
}

// SectionStatus represents the status of an onboarding section
type SectionStatus struct {
	Complete    bool      `bson:"complete" json:"complete"`
	LastUpdated time.Time `bson:"lastUpdated" json:"last_updated"`
}

// OnboardingStatus provides step completion flags, progress tracking, and audit timestamp
type OnboardingStatus struct {
	PercentComplete         int                      `bson:"percentComplete" json:"percent_complete"`
	NextStep                string                   `bson:"nextStep,omitempty" json:"next_step,omitempty"`
	Sections                map[string]SectionStatus `bson:"sections" json:"sections"`
	ConsentCompleted        bool                     `bson:"consentCompleted" json:"consent_completed"`
	DemographicsCompleted   bool                     `bson:"demographicsCompleted" json:"demographics_completed"`
	MedicalHistoryCompleted bool                     `bson:"medicalHistoryCompleted" json:"medical_history_completed"`
	LifestyleCompleted      bool                     `bson:"lifestyleCompleted" json:"lifestyle_completed"`
	LastUpdated             time.Time                `bson:"lastUpdated" json:"last_updated"`
}

// ConsentRecord represents a consent record (append-only)
type ConsentRecord struct {
	ID                      string                 `bson:"_id,omitempty"`
	UserID                  string                 `bson:"userId"`
	ConsentGiven            bool                   `bson:"consentGiven"`
	TermsVersion            string                 `bson:"termsVersion"`
	ConsentDate             time.Time              `bson:"consentDate"`
	DataSharingWithDoctors  bool                   `bson:"dataSharingWithDoctors"`
	DataSharingWithResearch bool                   `bson:"dataSharingWithResearch"`
	Signature               string                 `bson:"signature,omitempty"`
	CreatedAt               time.Time              `bson:"createdAt"`
	RevocationTimestamp     *time.Time             `bson:"revocationTimestamp,omitempty"`
	Metadata                map[string]interface{} `bson:"metadata,omitempty"`
}

// ConsentEnvelope captures granular consents and versioning + audit
type ConsentEnvelope struct {
	Baseline []ConsentItem `bson:"baseline" json:"baseline_consents"`
	Optional []ConsentItem `bson:"optional" json:"optional_consents"`
	// optional audit trail for changes over time
	AuditTrail  []ConsentItem `bson:"auditTrail,omitempty" json:"audit_trail,omitempty"`
	LastUpdated time.Time     `bson:"lastUpdated" json:"last_updated"`
}

type ConsentItem struct {
	Type      string    `bson:"type" json:"type"` // e.g., TERMS_OF_SERVICE, PRIVACY_POLICY, DATA_PROCESSING_FOR_SERVICE
	Granted   bool      `bson:"granted" json:"granted"`
	Version   string    `bson:"version,omitempty" json:"version,omitempty"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Reason    string    `bson:"reason,omitempty" json:"reason,omitempty"` // revocation reason
}

// Demographics stores standard identifiers with ISO/E.164 and encrypted PHI
type Demographics struct {
	FullName         encrypt.SecureString `bson:"fullName" json:"full_name"`        // PHI
	DateOfBirth      string               `bson:"dateOfBirth" json:"date_of_birth"` // ISO date (YYYY-MM-DD), can be encrypted if desired
	Gender           string               `bson:"gender" json:"gender"`
	Height           *Measure             `bson:"height,omitempty" json:"height,omitempty"` // unit: cm/in
	Weight           *Measure             `bson:"weight,omitempty" json:"weight,omitempty"` // unit: kg/lb
	BloodType        string               `bson:"bloodType,omitempty" json:"blood_type,omitempty"`
	Contact          *ContactInfo         `bson:"contact,omitempty" json:"contact,omitempty"`
	Location         *Location            `bson:"location,omitempty" json:"location,omitempty"`
	PrimaryLanguage  string               `bson:"primaryLanguage,omitempty" json:"primary_language,omitempty"` // ISO 639-1
	EmergencyContact *EmergencyContact    `bson:"emergencyContact,omitempty" json:"emergency_contact,omitempty"`
	LastUpdated      time.Time            `bson:"lastUpdated" json:"last_updated"`
}

type Measure struct {
	Value float64 `bson:"value" json:"value"`
	Unit  string  `bson:"unit" json:"unit"` // "cm","in","kg","lb"
}

type ContactInfo struct {
	Phone encrypt.SecureString `bson:"phone,omitempty" json:"phone,omitempty"` // E.164, PHI
	Email encrypt.SecureString `bson:"email,omitempty" json:"email,omitempty"` // PHI
}

type Location struct {
	Country string `bson:"country,omitempty" json:"country,omitempty"` // ISO 3166-1 alpha-2
	Region  string `bson:"region,omitempty" json:"region,omitempty"`
	City    string `bson:"city,omitempty" json:"city,omitempty"`
}

type EmergencyContact struct {
	Name     encrypt.SecureString `bson:"name" json:"name"` // PHI
	Relation string               `bson:"relation" json:"relation"`
	Phone    encrypt.SecureString `bson:"phone" json:"phone"` // E.164, PHI
}

// MedicalHistory uses structured arrays to enable standardized coding in future
type MedicalHistory struct {
	ExistingConditions []ConditionEntry `bson:"existingConditions,omitempty" json:"existing_conditions,omitempty"`
	PastSurgeries      []SurgeryEntry   `bson:"pastSurgeries,omitempty" json:"past_surgeries,omitempty"`
	CurrentMedications []Medication     `bson:"currentMedications,omitempty" json:"current_medications,omitempty"`
	Allergies          []Allergy        `bson:"allergies,omitempty" json:"allergies,omitempty"`
	FamilyHistory      []FamilyHistory  `bson:"familyHistory,omitempty" json:"family_history,omitempty"`
	LastUpdated        time.Time        `bson:"lastUpdated" json:"last_updated"`
}

type ConditionEntry struct {
	Condition     string `bson:"condition" json:"condition"`
	DiagnosedDate string `bson:"diagnosedDate,omitempty" json:"diagnosed_date,omitempty"`
	CodeSystem    string `bson:"codeSystem,omitempty" json:"code_system,omitempty"` // e.g., SNOMED_CT
	Code          string `bson:"code,omitempty" json:"code,omitempty"`
}

type SurgeryEntry struct {
	Surgery string `bson:"surgery" json:"surgery"`
	Date    string `bson:"date,omitempty" json:"date,omitempty"`
	Code    string `bson:"code,omitempty" json:"code,omitempty"`
}

type Medication struct {
	Name      string `bson:"name" json:"name"`
	Dosage    string `bson:"dosage,omitempty" json:"dosage,omitempty"`
	Frequency string `bson:"frequency,omitempty" json:"frequency,omitempty"`
	RxNorm    string `bson:"rxnorm,omitempty" json:"rxnorm,omitempty"`
}

type Allergy struct {
	Allergen string `bson:"allergen" json:"allergen"`
	Reaction string `bson:"reaction,omitempty" json:"reaction,omitempty"`
	Severity string `bson:"severity,omitempty" json:"severity,omitempty"` // Mild/Moderate/Severe
	SNOMED   string `bson:"snomed,omitempty" json:"snomed,omitempty"`
}

type FamilyHistory struct {
	Relation  string `bson:"relation" json:"relation"`
	Condition string `bson:"condition" json:"condition"`
	Code      string `bson:"code,omitempty" json:"code,omitempty"`
}

// Lifestyle uses enums for consistent analytics
type Lifestyle struct {
	DietaryHabits    *DietaryHabits    `bson:"dietaryHabits,omitempty" json:"dietary_habits,omitempty"`
	PhysicalActivity *PhysicalActivity `bson:"physicalActivity,omitempty" json:"physical_activity,omitempty"`
	SubstanceUse     *SubstanceUse     `bson:"substanceUse,omitempty" json:"substance_use,omitempty"`
	Sleep            *Sleep            `bson:"sleep,omitempty" json:"sleep,omitempty"`
	Stress           *Stress           `bson:"stress,omitempty" json:"stress,omitempty"`
	Occupation       string            `bson:"occupation,omitempty" json:"occupation,omitempty"`
	LastUpdated      time.Time         `bson:"lastUpdated" json:"last_updated"`
}

type DietaryHabits struct {
	DietType                string  `bson:"dietType,omitempty" json:"diet_type,omitempty"` // e.g., Omnivore, Vegetarian
	MealsPerDay             int     `bson:"mealsPerDay,omitempty" json:"meals_per_day,omitempty"`
	VegetableServingsPerDay int     `bson:"vegetableServingsPerDay,omitempty" json:"vegetable_intake_servings_per_day,omitempty"`
	DailyWaterIntakeLiters  float64 `bson:"dailyWaterIntakeLiters,omitempty" json:"daily_water_intake_liters,omitempty"`
}

type PhysicalActivity struct {
	FrequencyPerWeek int    `bson:"frequencyPerWeek,omitempty" json:"frequency_per_week,omitempty"`
	PrimaryType      string `bson:"primaryType,omitempty" json:"primary_type,omitempty"`
}

type SubstanceUse struct {
	SmokingStatus             string `bson:"smokingStatus,omitempty" json:"smoking_status,omitempty"` // Non-smoker, Former, Current
	AlcoholConsumptionPerWeek string `bson:"alcoholConsumptionPerWeek,omitempty" json:"alcohol_consumption_per_week,omitempty"`
}

type Sleep struct {
	AverageDurationHours float64 `bson:"averageDurationHours,omitempty" json:"average_duration_hours,omitempty"`
}

type Stress struct {
	Level string `bson:"level,omitempty" json:"level,omitempty"` // Low, Moderate, High
}
