package database

import (
	"fmt"
	"log"
	"os"

	"remedymate-backend/util"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client // will be initialized after ConnectMongo run

func ConnectMongo() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	mongoURI := os.Getenv("MONGO_URI")
	ctx, cancel := util.CreateContext()
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB connection error: ", err)
	}

	//Ping mongodb
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("MongoDB ping failed: ", err)
	}

	fmt.Println("MongoDB connected Seccussfully!")
	Client = client
}

// GetCollection returns a MongoDB collection by name
func GetCollection(collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "remedymate" // default database name
	}
	return Client.Database(dbName).Collection(collectionName)
}
