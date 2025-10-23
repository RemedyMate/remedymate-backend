package repository

import (
	"context"
	"time"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AppUserRepository struct {
	coll *mongo.Collection
}

func NewAppUserRepository() *AppUserRepository {
	c := database.GetCollection("app_users")
	// indexes
	ctx := context.Background()
	_, _ = c.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{"userId": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"status.lastUpdated": -1}},
		{Keys: bson.M{"updatedAt": -1}},
	})
	return &AppUserRepository{coll: c}
}

// GetOnboardingStatus returns the status (creates skeleton if not exists)
func (r *AppUserRepository) GetOnboardingStatus(ctx context.Context, userID string) (*entities.OnboardingStatus, error) {
	var doc entities.AppUser
	err := r.coll.FindOne(ctx, bson.M{"userId": userID}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		now := time.Now()
		doc = entities.AppUser{
			UserID: userID,
			Status: entities.OnboardingStatus{
				ConsentCompleted:        false,
				DemographicsCompleted:   false,
				MedicalHistoryCompleted: false,
				LifestyleCompleted:      false,
				LastUpdated:             now,
			},
			CreatedAt: now,
			UpdatedAt: now,
		}
		_, _ = r.coll.InsertOne(ctx, &doc)
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	return &doc.Status, nil
}

func (r *AppUserRepository) UpsertConsent(ctx context.Context, userID string, payload dto.ConsentDTO) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"userId": userID,
			"consent": bson.M{
				"baseline":    payload.Baseline,
				"optional":    payload.Optional,
				"lastUpdated": now,
				"auditTrail":  payload.AuditTrail, // optional
			},
			"status.consentCompleted": true,
			"status.lastUpdated":      now,
			"updatedAt":               now,
		},
		"$setOnInsert": bson.M{"createdAt": now},
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"userId": userID}, update, options.Update().SetUpsert(true))
	return err
}

func (r *AppUserRepository) UpdateDemographics(ctx context.Context, userID string, dem dto.DemographicsDTO) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"userId": userID,
			"demographics": bson.M{
				"fullName":         dem.FullName,
				"dateOfBirth":      dem.DateOfBirth,
				"gender":           dem.Gender,
				"height":           dem.Height,
				"weight":           dem.Weight,
				"bloodType":        dem.BloodType,
				"contact":          dem.Contact,
				"location":         dem.Location,
				"primaryLanguage":  dem.PrimaryLanguage,
				"emergencyContact": dem.EmergencyContact,
				"lastUpdated":      now,
			},
			"status.demographicsCompleted": true,
			"status.lastUpdated":           now,
			"updatedAt":                    now,
		},
		"$setOnInsert": bson.M{"createdAt": now},
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"userId": userID}, update, options.Update().SetUpsert(true))
	return err
}

func (r *AppUserRepository) UpdateMedicalHistory(ctx context.Context, userID string, mh dto.MedicalHistoryDTO) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"userId": userID,
			"medical_history": bson.M{
				"existingConditions": mh.ExistingConditions,
				"pastSurgeries":      mh.PastSurgeries,
				"currentMedications": mh.CurrentMedications,
				"allergies":          mh.Allergies,
				"familyHistory":      mh.FamilyHistory,
				"lastUpdated":        now,
			},
			"status.medicalHistoryCompleted": true,
			"status.lastUpdated":             now,
			"updatedAt":                      now,
		},
		"$setOnInsert": bson.M{"createdAt": now},
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"userId": userID}, update, options.Update().SetUpsert(true))
	return err
}

func (r *AppUserRepository) UpdateLifestyle(ctx context.Context, userID string, ls dto.LifestyleDTO) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"userId": userID,
			"lifestyle": bson.M{
				"dietaryHabits":    ls.DietaryHabits,
				"physicalActivity": ls.PhysicalActivity,
				"substanceUse":     ls.SubstanceUse,
				"sleep":            ls.Sleep,
				"stress":           ls.Stress,
				"occupation":       ls.Occupation,
				"lastUpdated":      now,
			},
			"status.lifestyleCompleted": true,
			"status.lastUpdated":        now,
			"updatedAt":                 now,
		},
		"$setOnInsert": bson.M{"createdAt": now},
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"userId": userID}, update, options.Update().SetUpsert(true))
	return err
}
