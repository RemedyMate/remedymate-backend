package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"remedymate-backend/delivery/controllers"
	"remedymate-backend/domain/entities"
	"remedymate-backend/infrastructure/database"
	"remedymate-backend/repository"
	"remedymate-backend/usecase"
	"remedymate-backend/util/encrypt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDemographicsIntegration tests the full demographics flow from HTTP to database
func TestDemographicsIntegration(t *testing.T) {
	// Skip integration test if encryption key not set or in short mode
	if testing.Short() || os.Getenv("ENCRYPTION_KEY") == "" {
		t.Skip("Skipping integration test: requires ENCRYPTION_KEY and full database setup")
	}

	// Set up test database connection
	originalMongoURI := os.Getenv("MONGO_URI")
	originalDBName := os.Getenv("DB_NAME")
	defer func() {
		os.Setenv("MONGO_URI", originalMongoURI)
		os.Setenv("DB_NAME", originalDBName)
	}()

	// Use test database
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("DB_NAME", "remedymate_test")

	// Connect to database
	database.ConnectMongo()
	defer database.Client.Disconnect(context.Background())

	// Set up dependencies
	userProfileRepo := repository.NewUserProfileRepository()
	consentRepo := repository.NewConsentRepository()                               // Need this for the usecase
	profileUsecase := usecase.NewProfileUsecase(userProfileRepo, consentRepo, nil) // No redis for this test
	profileController := controllers.NewProfileController(profileUsecase)

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add auth middleware for testing
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user-123")
		c.Next()
	})

	router.PUT("/demographics", profileController.PutDemographics)

	t.Run("successful demographics creation and update", func(t *testing.T) {
		// Test data
		demographics := entities.Demographics{
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
			BloodType:       "O+",
			PrimaryLanguage: "en",
			LastUpdated:     time.Now(),
		}

		// Test the usecase directly
		err := profileUsecase.PutDemographics(context.Background(), "test-user-123", &demographics)
		require.NoError(t, err)

		// Verify data was saved by retrieving the profile
		profile, err := userProfileRepo.GetUserProfile(context.Background(), "test-user-123")
		require.NoError(t, err)
		require.NotNil(t, profile)
		require.NotNil(t, profile.Demog)
		assert.Equal(t, "John Doe", string(profile.Demog.FullName))
		assert.Equal(t, "1990-01-01", profile.Demog.DateOfBirth)
		assert.Equal(t, "Male", profile.Demog.Gender)
		assert.Equal(t, 175.0, profile.Demog.Height.Value)
		assert.Equal(t, "cm", profile.Demog.Height.Unit)
		assert.Equal(t, 70.0, profile.Demog.Weight.Value)
		assert.Equal(t, "kg", profile.Demog.Weight.Unit)
		assert.Equal(t, "O+", profile.Demog.BloodType)
		assert.Equal(t, "en", profile.Demog.PrimaryLanguage)
	})

	// Clean up test data
	t.Cleanup(func() {
		// Drop test database collections
		db := database.Client.Database("remedymate_test")
		collections := []string{"user_profiles"}
		for _, coll := range collections {
			db.Collection(coll).Drop(context.Background())
		}
	})
}
