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

// TestPutLifestyleHandler tests the PUT /api/v1/user/profile/lifestyle endpoint
func TestPutLifestyleHandler(t *testing.T) {
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
			name:   "successful lifestyle update",
			userID: "user123",
			requestBody: entities.Lifestyle{
				DietaryHabits: &entities.DietaryHabits{
					DietType:                "Mediterranean",
					MealsPerDay:             3,
					VegetableServingsPerDay: 5,
					DailyWaterIntakeLiters:  2.5,
				},
				PhysicalActivity: &entities.PhysicalActivity{
					FrequencyPerWeek: 4,
					PrimaryType:      "Running",
				},
				SubstanceUse: &entities.SubstanceUse{
					SmokingStatus:             "Non-smoker",
					AlcoholConsumptionPerWeek: "1-2 drinks",
				},
				Sleep: &entities.Sleep{
					AverageDurationHours: 7.5,
				},
				Stress: &entities.Stress{
					Level: "Moderate",
				},
				Occupation:  "Software Engineer",
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutLifestyle", mock.Anything, "user123", mock.AnythingOfType("*entities.Lifestyle")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Lifestyle updated successfully"}`,
		},
		{
			name:   "minimal lifestyle update",
			userID: "user123",
			requestBody: entities.Lifestyle{
				Occupation:  "Teacher",
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutLifestyle", mock.Anything, "user123", mock.AnythingOfType("*entities.Lifestyle")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Lifestyle updated successfully"}`,
		},
		{
			name:   "lifestyle with only dietary habits",
			userID: "user123",
			requestBody: entities.Lifestyle{
				DietaryHabits: &entities.DietaryHabits{
					DietType:                "Vegetarian",
					MealsPerDay:             2,
					VegetableServingsPerDay: 8,
					DailyWaterIntakeLiters:  3.0,
				},
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutLifestyle", mock.Anything, "user123", mock.AnythingOfType("*entities.Lifestyle")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Lifestyle updated successfully"}`,
		},
		{
			name:   "lifestyle with only physical activity",
			userID: "user123",
			requestBody: entities.Lifestyle{
				PhysicalActivity: &entities.PhysicalActivity{
					FrequencyPerWeek: 3,
					PrimaryType:      "Swimming",
				},
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutLifestyle", mock.Anything, "user123", mock.AnythingOfType("*entities.Lifestyle")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Lifestyle updated successfully"}`,
		},
		{
			name:   "missing user ID",
			userID: "",
			requestBody: entities.Lifestyle{
				Occupation:  "Doctor",
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
			requestBody: entities.Lifestyle{
				Occupation:  "Nurse",
				LastUpdated: time.Now(),
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PutLifestyle", mock.Anything, "user123", mock.AnythingOfType("*entities.Lifestyle")).Return(assert.AnError)
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
			router.PUT("/lifestyle", controller.PutLifestyle)

			// Create request
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}
			req, _ := http.NewRequest(http.MethodPut, "/lifestyle", bytes.NewBuffer(body))
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
