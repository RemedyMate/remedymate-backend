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
	"remedymate-backend/util/encrypt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestPutDemographicsHandler tests the PUT /api/v1/user/profile/demographics endpoint
func TestPutDemographicsHandler(t *testing.T) {
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
			name:   "successful demographics update",
			userID: "user123",
			requestBody: entities.Demographics{
				FullName:    encrypt.SecureString("John Doe"),
				DateOfBirth: "1990-01-01",
				Gender:      "Male",
				Height: &entities.Measure{
					Value: 175.0,
					Unit:  "cm",
				},
				Weight: &entities.Measure{
					Value: 70.0,
					Unit:  "kg",
				},
				BloodType: "O+",
				Contact: &entities.ContactInfo{
					Phone: encrypt.SecureString("+1234567890"),
					Email: encrypt.SecureString("john.doe@example.com"),
				},
				Location: &entities.Location{
					Country: "US",
					Region:  "CA",
					City:    "San Francisco",
				},
				PrimaryLanguage: "en",
				EmergencyContact: &entities.EmergencyContact{
					Name:     encrypt.SecureString("Jane Doe"),
					Relation: "Spouse",
					Phone:    encrypt.SecureString("+0987654321"),
				},
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutDemographics", mock.Anything, "user123", mock.AnythingOfType("*entities.Demographics")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Demographics updated successfully"}`,
		},
		{
			name:   "minimal demographics update",
			userID: "user123",
			requestBody: entities.Demographics{
				FullName:    encrypt.SecureString("John Doe"),
				DateOfBirth: "1990-01-01",
				Gender:      "Male",
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutDemographics", mock.Anything, "user123", mock.AnythingOfType("*entities.Demographics")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Demographics updated successfully"}`,
		},
		{
			name:   "missing user ID",
			userID: "",
			requestBody: entities.Demographics{
				FullName:    encrypt.SecureString("John Doe"),
				DateOfBirth: "1990-01-01",
				Gender:      "Male",
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
			requestBody: entities.Demographics{
				FullName:    encrypt.SecureString("John Doe"),
				DateOfBirth: "1990-01-01",
				Gender:      "Male",
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutDemographics", mock.Anything, "user123", mock.AnythingOfType("*entities.Demographics")).Return(assert.AnError)
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
			router.PUT("/demographics", controller.PutDemographics)

			// Create request
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}
			req, _ := http.NewRequest(http.MethodPut, "/demographics", bytes.NewBuffer(body))
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
