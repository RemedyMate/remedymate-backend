package domain

import "time"

// UserProfile represents the complete user profile with onboarding data
type UserProfile struct {
	ID               string            `bson:"_id,omitempty"`
	UserID           string            `bson:"userId"`
	Demographics     *Demographics     `bson:"demographics,omitempty"`
	MedicalHistory   *MedicalHistory   `bson:"medicalHistory,omitempty"`
	Lifestyle        *Lifestyle        `bson:"lifestyle,omitempty"`
	OnboardingStatus *OnboardingStatus `bson:"onboardingStatus,omitempty"`
	CreatedAt        time.Time         `bson:"createdAt"`
	UpdatedAt        time.Time         `bson:"updatedAt"`
}

// Demographics represents user demographic information
type Demographics struct {
	FullName                 string    `bson:"fullName"`
	Age                      int       `bson:"age"`
	Gender                   string    `bson:"gender,omitempty"`
	Height                   float64   `bson:"height,omitempty"`
	Weight                   float64   `bson:"weight,omitempty"`
	BloodType                string    `bson:"bloodType,omitempty"`
	Email                    string    `bson:"email,omitempty"`
	Phone                    string    `bson:"phone,omitempty"`
	Region                   string    `bson:"region,omitempty"`
	City                     string    `bson:"city,omitempty"`
	PrimaryLanguage          string    `bson:"primaryLanguage,omitempty"`
	EmergencyContactName     string    `bson:"emergencyContactName,omitempty"`
	EmergencyContactRelation string    `bson:"emergencyContactRelation,omitempty"`
	EmergencyContactPhone    string    `bson:"emergencyContactPhone,omitempty"`
	LastUpdated              time.Time `bson:"lastUpdated"`
}

// MedicalHistory represents user's medical history
type MedicalHistory struct {
	ExistingConditions     []string  `bson:"existingConditions,omitempty"`
	PastSurgeries          []string  `bson:"pastSurgeries,omitempty"`
	CurrentMedications     []string  `bson:"currentMedications,omitempty"`
	Allergies              []string  `bson:"allergies,omitempty"`
	FamilyHistory          []string  `bson:"familyHistory,omitempty"`
	ChronicPain            bool      `bson:"chronicPain"`
	RecentHospitalizations []string  `bson:"recentHospitalizations,omitempty"`
	LastUpdated            time.Time `bson:"lastUpdated"`
}

// Lifestyle represents user's lifestyle data
type Lifestyle struct {
	DietType                 string    `bson:"dietType,omitempty"`
	MealsPerDay              int       `bson:"mealsPerDay,omitempty"`
	VegetableIntake          string    `bson:"vegetableIntake,omitempty"`
	WaterIntakeLiters        float64   `bson:"waterIntakeLiters,omitempty"`
	ExerciseFrequencyPerWeek int       `bson:"exerciseFrequencyPerWeek,omitempty"`
	ExerciseType             string    `bson:"exerciseType,omitempty"`
	SmokingStatus            string    `bson:"smokingStatus,omitempty"`
	AlcoholConsumption       string    `bson:"alcoholConsumption,omitempty"`
	SleepDurationHours       float64   `bson:"sleepDurationHours,omitempty"`
	StressLevel              string    `bson:"stressLevel,omitempty"`
	Occupation               string    `bson:"occupation,omitempty"`
	ScreenTimeHours          float64   `bson:"screenTimeHours,omitempty"`
	LastUpdated              time.Time `bson:"lastUpdated"`
}

// OnboardingStatus represents the onboarding progress
type OnboardingStatus struct {
	PercentComplete int                      `bson:"percentComplete"`
	NextStep        string                   `bson:"nextStep,omitempty"`
	Sections        map[string]SectionStatus `bson:"sections"`
}

// SectionStatus represents the status of an onboarding section
type SectionStatus struct {
	Complete    bool      `bson:"complete"`
	LastUpdated time.Time `bson:"lastUpdated"`
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
