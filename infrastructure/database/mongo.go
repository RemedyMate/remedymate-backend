package database

import (
	"fmt"
	"log"
	"os"

	"github.com/RemedyMate/remedymate-backend/util"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client

func ConnectMongo() *mongo.Client{
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
	
	return client
}

