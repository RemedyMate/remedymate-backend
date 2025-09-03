package main

import (
	"context"
	"encoding/json"
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
	data, err := os.ReadFile("data/approved_blocks.json")
	if err != nil {
		log.Fatal("❌ Failed to read seed file:", err)
	}

	// 3️⃣ Temporary struct to parse IDs as strings and dates as strings
	type rawTopic struct {
		ID            string                `json:"_id,omitempty"`
		TopicKey      string                `json:"topic_key"`
		NameEn        string                `json:"name_en,omitempty"`
		NameAm        string                `json:"name_am,omitempty"`
		DescriptionEn string                `json:"description_en,omitempty"`
		DescriptionAm string                `json:"description_am,omitempty"`
		Status        string                `json:"status,omitempty"`
		Translations  entities.Translations `json:"translations"`
		Version       int                   `json:"version,omitempty"`
		CreatedAt     string                `json:"created_at,omitempty"`
		UpdatedAt     string                `json:"updated_at,omitempty"`
		CreatedBy     string                `json:"created_by,omitempty"`
		UpdatedBy     string                `json:"updated_by,omitempty"`
	}

	var rawTopics []rawTopic
	if err := json.Unmarshal(data, &rawTopics); err != nil {
		log.Fatal("❌ Failed to parse JSON:", err)
	}

	// 4️⃣ Get collection and prepare context
	collection := database.GetCollection("health_topics")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Default values for missing fields
	defaultStatus := "active"
	defaultVersion := 1
	defaultCreatedBy := primitive.NewObjectID() // Generate a default user ID
	defaultUpdatedBy := defaultCreatedBy
	currentTime := time.Now()

	// 5️⃣ Process each topic
	var topics []entities.HealthTopic
	for _, raw := range rawTopics {
		// Generate ID if not provided
		var id primitive.ObjectID
		if raw.ID == "" {
			id = primitive.NewObjectID()
			log.Println("ℹ️  Generated new ID for topic", raw.TopicKey+":", id.Hex())
		} else {
			var err error
			id, err = primitive.ObjectIDFromHex(raw.ID)
			if err != nil {
				log.Fatal("❌ Invalid _id for topic", raw.TopicKey+":", err)
			}
		}

		// Handle createdBy
		var createdBy primitive.ObjectID
		if raw.CreatedBy == "" {
			createdBy = defaultCreatedBy
		} else {
			var err error
			createdBy, err = primitive.ObjectIDFromHex(raw.CreatedBy)
			if err != nil {
				log.Fatal("❌ Invalid created_by for topic", raw.TopicKey+":", err)
			}
		}

		// Handle updatedBy
		var updatedBy primitive.ObjectID
		if raw.UpdatedBy == "" {
			updatedBy = defaultUpdatedBy
		} else {
			var err error
			updatedBy, err = primitive.ObjectIDFromHex(raw.UpdatedBy)
			if err != nil {
				log.Fatal("❌ Invalid updated_by for topic", raw.TopicKey+":", err)
			}
		}

		// Handle timestamps
		var createdAt, updatedAt time.Time
		if raw.CreatedAt == "" {
			createdAt = currentTime
		} else {
			var err error
			createdAt, err = time.Parse(time.RFC3339, raw.CreatedAt)
			if err != nil {
				log.Fatal("❌ Invalid created_at for topic", raw.TopicKey+":", err)
			}
		}

		if raw.UpdatedAt == "" {
			updatedAt = currentTime
		} else {
			var err error
			updatedAt, err = time.Parse(time.RFC3339, raw.UpdatedAt)
			if err != nil {
				log.Fatal("❌ Invalid updated_at for topic", raw.TopicKey+":", err)
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

		topic := entities.HealthTopic{
			ID:            id,
			TopicKey:      raw.TopicKey,
			NameEn:        raw.NameEn,
			NameAm:        raw.NameAm,
			DescriptionEn: raw.DescriptionEn,
			DescriptionAm: raw.DescriptionAm,
			Status:        status,
			Translations:  raw.Translations,
			Version:       version,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
			CreatedBy:     createdBy,
			UpdatedBy:     updatedBy,
		}

		topics = append(topics, topic)
	}

	// 6️⃣ Insert all documents into MongoDB
	var documents []interface{}
	for _, topic := range topics {
		documents = append(documents, topic)
	}

	result, err := collection.InsertMany(ctx, documents)
	if err != nil {
		log.Fatal("❌ Failed to insert documents:", err)
	}

	log.Printf("✅ Successfully inserted %d documents into database", len(result.InsertedIDs))
	log.Println("Inserted IDs:", result.InsertedIDs)
}
