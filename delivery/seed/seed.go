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
	data, err := os.ReadFile("delivery/seed/seed.json")
	if err != nil {
		log.Fatal("❌ Failed to read seed file:", err)
	}

	// 3️⃣ Temporary struct to parse IDs as strings and dates as strings
	type rawTopic struct {
		ID            string                `json:"_id"`
		TopicKey      string                `json:"topic_key"`
		NameEn        string                `json:"name_en"`
		NameAm        string                `json:"name_am"`
		DescriptionEn string                `json:"description_en"`
		DescriptionAm string                `json:"description_am"`
		Status        string                `json:"status"`
		Translations  entities.Translations `json:"translations"`
		Version       int                   `json:"version"`
		CreatedAt     string                `json:"created_at"`
		UpdatedAt     string                `json:"updated_at"`
		CreatedBy     string                `json:"created_by"`
		UpdatedBy     string                `json:"updated_by"`
	}

	var raw rawTopic
	if err := json.Unmarshal(data, &raw); err != nil {
		log.Fatal("❌ Failed to parse JSON:", err)
	}

	createdBy, err := primitive.ObjectIDFromHex(raw.CreatedBy)
	if err != nil {
		log.Fatal("❌ Invalid created_by:", err)
	}
	updatedBy, err := primitive.ObjectIDFromHex(raw.UpdatedBy)
	if err != nil {
		log.Fatal("❌ Invalid updated_by:", err)
	}

	createdAt, err := time.Parse(time.RFC3339, raw.CreatedAt)
	if err != nil {
		log.Fatal("❌ Invalid created_at:", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, raw.UpdatedAt)
	if err != nil {
		log.Fatal("❌ Invalid updated_at:", err)
	}

	topic := entities.HealthTopic{
		TopicKey:      raw.TopicKey,
		NameEn:        raw.NameEn,
		NameAm:        raw.NameAm,
		DescriptionEn: raw.DescriptionEn,
		DescriptionAm: raw.DescriptionAm,
		Status:        raw.Status,
		Translations:  raw.Translations,
		Version:       raw.Version,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		CreatedBy:     createdBy,
		UpdatedBy:     updatedBy,
	}

	// 5️⃣ Insert into MongoDB
	collection := database.GetCollection("health_topics")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, topic)
	if err != nil {
		log.Fatal("❌ Failed to insert document:", err)
	}

	log.Println("✅ Database seeded successfully")
}
