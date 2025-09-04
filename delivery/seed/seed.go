package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"remedymate-backend/domain/entities"
	"remedymate-backend/infrastructure/database"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// 1️⃣ Load .env and connect to MongoDB
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not found")
	}
	database.ConnectMongo()

	// 2️⃣ Load JSON file
	data, err := os.ReadFile("data/approved_block.json")
	if err != nil {
		log.Fatal("❌ Failed to read seed file:", err)
	}

	// MongoDB ObjectID wrapper
	type ObjectIDWrapper struct {
		OID string `json:"$oid"`
	}

	// MongoDB Date wrapper
	type DateWrapper struct {
		Date string `json:"$date"`
	}

	// 3️⃣ Temporary struct to parse MongoDB extended JSON format
	type rawTopic struct {
		ID            ObjectIDWrapper       `json:"_id,omitempty"`
		TopicKey      string                `json:"topic_key"`
		NameEn        string                `json:"name_en,omitempty"`
		NameAm        string                `json:"name_am,omitempty"`
		DescriptionEn string                `json:"description_en,omitempty"`
		DescriptionAm string                `json:"description_am,omitempty"`
		Status        string                `json:"status,omitempty"`
		Translations  entities.Translations `json:"translations"`
		Version       int                   `json:"version,omitempty"`
		CreatedAt     DateWrapper           `json:"created_at,omitempty"`
		UpdatedAt     DateWrapper           `json:"updated_at,omitempty"`
		CreatedBy     ObjectIDWrapper       `json:"created_by,omitempty"`
		UpdatedBy     ObjectIDWrapper       `json:"updated_by,omitempty"`
	}

	var rawTopics []rawTopic
	if err := json.Unmarshal(data, &rawTopics); err != nil {
		log.Fatal("❌ Failed to parse JSON:", err)
	}

	// 4️⃣ Get collection and prepare context
	collection := database.GetCollection("health_topics")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Use 30s timeout to accommodate potentially large batch insertions
	defer cancel()

	// Default values for missing fields
	defaultStatus := "active"
	defaultVersion := 1
	defaultCreatedBy := primitive.NewObjectID() // Generate a default user ID
	defaultUpdatedBy := defaultCreatedBy
	currentTime := time.Now()

	// 5️⃣ Process each topic and prepare for insertion
	var documents []interface{}
	successCount := 0

	for _, raw := range rawTopics {
		// Handle ID
		var id primitive.ObjectID
		if raw.ID.OID == "" {
			id = primitive.NewObjectID()
			log.Printf("ℹ️  Generated new ID for topic %s: %s", raw.TopicKey, id.Hex())
		} else {
			var err error
			id, err = primitive.ObjectIDFromHex(raw.ID.OID)
			if err != nil {
				log.Printf("❌ Invalid _id for topic %s: %v", raw.TopicKey, err)
				continue // Skip this topic but continue with others
			}
		}

		// Handle createdBy
		var createdBy primitive.ObjectID
		if raw.CreatedBy.OID == "" {
			createdBy = defaultCreatedBy
		} else {
			var err error
			createdBy, err = primitive.ObjectIDFromHex(raw.CreatedBy.OID)
			if err != nil {
				log.Printf("❌ Invalid created_by for topic %s: %v", raw.TopicKey, err)
				continue
			}
		}

		// Handle updatedBy
		var updatedBy primitive.ObjectID
		if raw.UpdatedBy.OID == "" {
			updatedBy = defaultUpdatedBy
		} else {
			var err error
			updatedBy, err = primitive.ObjectIDFromHex(raw.UpdatedBy.OID)
			if err != nil {
				log.Printf("❌ Invalid updated_by for topic %s: %v", raw.TopicKey, err)
				continue
			}
		}

		// Handle timestamps
		var createdAt, updatedAt time.Time
		if raw.CreatedAt.Date == "" {
			createdAt = currentTime
		} else {
			var err error
			createdAt, err = time.Parse(time.RFC3339, raw.CreatedAt.Date)
			if err != nil {
				log.Printf("❌ Invalid created_at for topic %s: %v", raw.TopicKey, err)
				continue
			}
		}

		if raw.UpdatedAt.Date == "" {
			updatedAt = currentTime
		} else {
			var err error
			updatedAt, err = time.Parse(time.RFC3339, raw.UpdatedAt.Date)
			if err != nil {
				log.Printf("❌ Invalid updated_at for topic %s: %v\n", raw.TopicKey, err)
				continue
			}
		}

		// Set default values for other optional fields
		status := raw.Status
		if status == "" {
			status = defaultStatus
		}

		version := raw.Version
		if version == 0 {
			version = defaultVersion
		}

		// Set default names and descriptions based on topic key if not provided
		nameEn := raw.NameEn
		if nameEn == "" {
			// You can add logic to generate default names based on topic_key
			nameEn = fmt.Sprintf("Default English Name for %s", raw.TopicKey)
		}

		nameAm := raw.NameAm

		descriptionEn := raw.DescriptionEn
		if descriptionEn == "" {
			descriptionEn = "Default English description for " + raw.TopicKey
		}

		descriptionAm := raw.DescriptionAm

		topic := entities.HealthTopic{
			ID:            id,
			TopicKey:      raw.TopicKey,
			NameEn:        nameEn,
			NameAm:        nameAm,
			DescriptionEn: descriptionEn,
			DescriptionAm: descriptionAm,
			Status:        status,
			Translations:  raw.Translations,
			Version:       version,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
			CreatedBy:     createdBy,
			UpdatedBy:     updatedBy,
		}

		documents = append(documents, topic)
		successCount++
	}

	if len(documents) == 0 {
		log.Fatal("❌ No valid documents to insert")
	}

	// 6️⃣ Insert all documents into MongoDB
	result, err := collection.InsertMany(ctx, documents)
	if err != nil {
		log.Fatal("❌ Failed to insert documents:", err)
	}

	log.Printf("✅ Successfully inserted %d out of %d documents into database", len(result.InsertedIDs), len(rawTopics))
	log.Println("Inserted IDs:", result.InsertedIDs)
}
