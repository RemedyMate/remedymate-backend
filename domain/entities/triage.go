package entities

// TriageLevel represents the severity level of symptoms
type TriageLevel string

const (
	TriageLevelGreen  TriageLevel = "GREEN"  // Likely mild
	TriageLevelYellow TriageLevel = "YELLOW" // Monitor closely
	TriageLevelRed    TriageLevel = "RED"    // Seek urgent care now
)

// TriageResult represents the result of symptom triage
type TriageResult struct {
	Level    TriageLevel `json:"level" bson:"level"`
	RedFlags []string    `json:"red_flags" bson:"red_flags"`
	Message  string      `json:"message" bson:"message"`
}

// SymptomInput represents user input for symptoms
type SymptomInput struct {
	Text     string `json:"text" bson:"text"`
	Language string `json:"language" bson:"language"` // "en" or "am"
}
