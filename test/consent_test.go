package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"remedymate-backend/delivery/controllers"
	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// testAuthMiddleware simulates authentication middleware for testing
func testAuthMiddleware(userID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if userID != "" {
			c.Set("userID", userID)
		}
		c.Next()
	}
}

// MockProfileUsecase is a mock implementation of IProfileUsecase for testing
type MockProfileUsecase struct {
	mock.Mock
}

func (m *MockProfileUsecase) GetOnboardingStatus(ctx context.Context, userID string) (*entities.OnboardingStatus, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*entities.OnboardingStatus), args.Error(1)
}

func (m *MockProfileUsecase) PostConsent(ctx context.Context, userID string, envelope *entities.ConsentEnvelope) error {
	args := m.Called(ctx, userID, envelope)
	return args.Error(0)
}

func (m *MockProfileUsecase) PatchConsent(ctx context.Context, userID string, consentType, reason string) error {
	args := m.Called(ctx, userID, consentType, reason)
	return args.Error(0)
}

func (m *MockProfileUsecase) PutDemographics(ctx context.Context, userID string, demographics *entities.Demographics) error {
	args := m.Called(ctx, userID, demographics)
	return args.Error(0)
}

func (m *MockProfileUsecase) PutMedicalHistory(ctx context.Context, userID string, history *entities.MedicalHistory) error {
	args := m.Called(ctx, userID, history)
	return args.Error(0)
}

func (m *MockProfileUsecase) PutLifestyle(ctx context.Context, userID string, lifestyle *entities.Lifestyle) error {
	args := m.Called(ctx, userID, lifestyle)
	return args.Error(0)
}

// TestPostConsentHandler tests the POST /api/v1/user/profile/consent endpoint
func TestPostConsentHandler(t *testing.T) {
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
			name:   "successful consent recording",
			userID: "user123",
			requestBody: dto.ConsentPostRequest{
				Consents: []dto.ConsentItemDTO{
					{
						Type:    "TERMS_OF_SERVICE",
						Granted: true,
						Version: "1.0",
					},
				},
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PostConsent", mock.Anything, "user123", mock.AnythingOfType("*entities.ConsentEnvelope")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message":"Consent recorded successfully"}`,
		},
		{
			name:   "missing user ID",
			userID: "",
			requestBody: dto.ConsentPostRequest{
				Consents: []dto.ConsentItemDTO{
					{
						Type:    "TERMS_OF_SERVICE",
						Granted: true,
						Version: "1.0",
					},
				},
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
			requestBody: dto.ConsentPostRequest{
				Consents: []dto.ConsentItemDTO{
					{
						Type:    "TERMS_OF_SERVICE",
						Granted: true,
						Version: "1.0",
					},
				},
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PostConsent", mock.Anything, "user123", mock.AnythingOfType("*entities.ConsentEnvelope")).Return(assert.AnError)
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
			router.POST("/consent", controller.PostConsent)

			// Create request
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}
			req, _ := http.NewRequest(http.MethodPost, "/consent", bytes.NewBuffer(body))
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

// TestPatchConsentHandler tests the PATCH /api/v1/user/profile/consent endpoint
func TestPatchConsentHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		requestBody    interface{}
		setupMocks     func(*MockProfileUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful consent revocation",
			userID: "user123",
			requestBody: dto.ConsentPatchRequest{
				Consents: []dto.ConsentRevocationDTO{
					{
						Type:             "TERMS_OF_SERVICE",
						Granted:          false,
						RevocationReason: "User requested withdrawal",
					},
				},
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PatchConsent", mock.Anything, "user123", "TERMS_OF_SERVICE", "User requested withdrawal").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Consents updated successfully"}`,
		},
		{
			name:   "multiple revocations",
			userID: "user123",
			requestBody: dto.ConsentPatchRequest{
				Consents: []dto.ConsentRevocationDTO{
					{
						Type:             "TERMS_OF_SERVICE",
						Granted:          false,
						RevocationReason: "Changed mind",
					},
					{
						Type:             "DATA_PROCESSING",
						Granted:          false,
						RevocationReason: "Privacy concerns",
					},
				},
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PatchConsent", mock.Anything, "user123", "TERMS_OF_SERVICE", "Changed mind").Return(nil)
				m.On("PatchConsent", mock.Anything, "user123", "DATA_PROCESSING", "Privacy concerns").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Consents updated successfully"}`,
		},
		{
			name:   "missing user ID",
			userID: "",
			requestBody: dto.ConsentPatchRequest{
				Consents: []dto.ConsentRevocationDTO{
					{
						Type:             "TERMS_OF_SERVICE",
						Granted:          false,
						RevocationReason: "Test",
					},
				},
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
			name:   "usecase error on first revocation",
			userID: "user123",
			requestBody: dto.ConsentPatchRequest{
				Consents: []dto.ConsentRevocationDTO{
					{
						Type:             "TERMS_OF_SERVICE",
						Granted:          false,
						RevocationReason: "Error test",
					},
				},
			},
			setupMocks: func(m *MockProfileUsecase) {
				m.On("PatchConsent", mock.Anything, "user123", "TERMS_OF_SERVICE", "Error test").Return(assert.AnError)
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
			router.PATCH("/consent", controller.PatchConsent)

			// Create request
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}
			req, _ := http.NewRequest(http.MethodPatch, "/consent", bytes.NewBuffer(body))
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
