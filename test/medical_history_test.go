package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"remedymate-backend/delivery/controllers"
	"remedymate-backend/domain/entities"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestPutMedicalHistoryHandler tests the PUT /api/v1/user/profile/medical-history endpoint
func TestPutMedicalHistoryHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		requestBody    interface{} // Changed to interface{} to allow invalid JSON
		setupMocks     func(*MockProfileUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful medical history update",
			userID: "user123",
			requestBody: entities.MedicalHistory{
				ExistingConditions: []entities.ConditionEntry{
					{
						Condition:     "Hypertension",
						DiagnosedDate: "2020-01-01",
						CodeSystem:    "SNOMED_CT",
						Code:          "38341003",
					},
				},
				PastSurgeries: []entities.SurgeryEntry{
					{
						Surgery: "Appendectomy",
						Date:    "2019-05-15",
						Code:    "80146002",
					},
				},
				CurrentMedications: []entities.Medication{
					{
						Name:      "Lisinopril",
						Dosage:    "10mg",
						Frequency: "Once daily",
						RxNorm:    "314076",
					},
				},
				Allergies: []entities.Allergy{
					{
						Allergen: "Penicillin",
						Reaction: "Rash",
						Severity: "Moderate",
						SNOMED:   "91936005",
					},
				},
				FamilyHistory: []entities.FamilyHistory{
					{
						Relation:  "Father",
						Condition: "Diabetes",
						Code:      "73211009",
					},
				},
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutMedicalHistory", mock.Anything, "user123", mock.AnythingOfType("*entities.MedicalHistory")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Medical history updated successfully"}`,
		},
		{
			name:   "minimal medical history update",
			userID: "user123",
			requestBody: entities.MedicalHistory{
				ExistingConditions: []entities.ConditionEntry{
					{
						Condition: "Asthma",
					},
				},
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutMedicalHistory", mock.Anything, "user123", mock.AnythingOfType("*entities.MedicalHistory")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Medical history updated successfully"}`,
		},
		{
			name:   "empty medical history update",
			userID: "user123",
			requestBody: entities.MedicalHistory{
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutMedicalHistory", mock.Anything, "user123", mock.AnythingOfType("*entities.MedicalHistory")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Medical history updated successfully"}`,
		},
		{
			name:   "missing user ID",
			userID: "",
			requestBody: entities.MedicalHistory{
				ExistingConditions: []entities.ConditionEntry{
					{
						Condition: "Diabetes",
					},
				},
				LastUpdated: time.Now(),
			},
			setupMocks:     func(m *MockProfileUsecase) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"User not authenticated"}`,
		},
		{
			name:           "invalid request body",
			userID:         "user123",
			requestBody:    "invalid json",
			setupMocks:     func(m *MockProfileUsecase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request body"}`,
		},
		{
			name:   "usecase error",
			userID: "user123",
			requestBody: entities.MedicalHistory{
				ExistingConditions: []entities.ConditionEntry{
					{
						Condition: "Migraine",
					},
				},
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutMedicalHistory", mock.Anything, "user123", mock.AnythingOfType("*entities.MedicalHistory")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"assert.AnError general error for testing"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock usecase
			mockUsecase := &MockProfileUsecase{}
			tt.setupMocks(mockUsecase)

			// Create controller
			controller := controllers.NewProfileController(mockUsecase)

			// Create gin router
			router := gin.New()
			router.Use(testAuthMiddleware(tt.userID))
			router.PUT("/medical-history", controller.PutMedicalHistory)

			// Create request
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}
			req, _ := http.NewRequest(http.MethodPut, "/medical-history", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockUsecase.AssertExpectations(t)
		})
	}
}
