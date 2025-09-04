package entities

import (
	"time"
)

// ConversationStatus represents the status of a conversation
type ConversationStatus string

const (
	ConversationStatusActive   ConversationStatus = "ACTIVE"
	ConversationStatusComplete ConversationStatus = "COMPLETE"
	ConversationStatusExpired  ConversationStatus = "EXPIRED"
)

// Conversation represents a conversation session
type Conversation struct {
	ID          string             `json:"id" bson:"_id"`
	UserID      string             `json:"user_id,omitempty" bson:"user_id,omitempty"` // Optional, for unauthenticated users
	Symptom     string             `json:"symptom" bson:"symptom"`
	Language    string             `json:"language" bson:"language"`
	Status      ConversationStatus `json:"status" bson:"status"`
	Questions   []Question         `json:"questions" bson:"questions"`
	Answers     []Answer           `json:"answers" bson:"answers"`
	CurrentStep int                `json:"current_step" bson:"current_step"`
	TotalSteps  int                `json:"total_steps" bson:"total_steps"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	CompletedAt *time.Time         `json:"completed_at" bson:"completed_at,omitempty"`
	FinalReport *HealthReport      `json:"final_report" bson:"final_report,omitempty"`
}

// Question represents a follow-up question in the conversation
type Question struct {
	ID       int    `json:"id" bson:"id"`
	Text     string `json:"text" bson:"text"`
	Type     string `json:"type" bson:"type"` // "duration", "location", "severity", "history", "triggers"
	Required bool   `json:"required" bson:"required"`
}

// Answer represents a user's answer to a question
type Answer struct {
	QuestionID int       `json:"question_id" bson:"question_id"`
	Text       string    `json:"text" bson:"text"`
	IsValid    bool      `json:"is_valid" bson:"is_valid"`
	Feedback   string    `json:"feedback" bson:"feedback,omitempty"`
	AnsweredAt time.Time `json:"answered_at" bson:"answered_at"`
}

// Remedy represents remedy information from the GetRemedy usecase
type Remedy struct {
	Triage        TriageResult  `json:"triage" bson:"triage"`
	SelfCare      []string      `json:"self_care" bson:"self_care"`
	OTCCategories []OTCCategory `json:"otc_categories,omitempty" bson:"otc_categories,omitempty"`
	SeekCareIf    []string      `json:"seek_care_if" bson:"seek_care_if"`
	Disclaimer    string        `json:"disclaimer" bson:"disclaimer"`
	TopicKey      string        `json:"topic_key,omitempty" bson:"topic_key,omitempty"`
	Language      string        `json:"language,omitempty" bson:"language,omitempty"`
}

// HealthReport represents the final structured health report
type HealthReport struct {
	Symptom            string    `json:"symptom" bson:"symptom"`
	Duration           string    `json:"duration" bson:"duration"`
	Location           string    `json:"location" bson:"location"`
	Severity           string    `json:"severity" bson:"severity"`
	AssociatedSymptoms []string  `json:"associated_symptoms" bson:"associated_symptoms"`
	MedicalHistory     string    `json:"medical_history" bson:"medical_history"`
	Triggers           string    `json:"triggers" bson:"triggers"`
	PossibleConditions []string  `json:"possible_conditions" bson:"possible_conditions"`
	Recommendations    []string  `json:"recommendations" bson:"recommendations"`
	UrgencyLevel       string    `json:"urgency_level" bson:"urgency_level"`
	GeneratedAt        time.Time `json:"generated_at" bson:"generated_at"`
	Remedy             *Remedy   `json:"remedy,omitempty" bson:"remedy,omitempty"`
}
