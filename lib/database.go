package lib

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// "github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *mongo.Client
	DB 	   *mongo.Database
}

func InitializeMongoDB() error {
	
	// if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
    //     if err := godotenv.Load(); err != nil {
    //         log.Println("No .env file found, relying on Railway variables")
    //     }
    // }
	mongoURI := os.Getenv("DATABASE_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err) 
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		client.Disconnect(context.Background())
        return fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	log.Println("database connected")
	
	Client = client
	DB = client.Database("url_shorten")
	

	return nil
}

var (
    DB     *mongo.Database 
    Client *mongo.Client
)

func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return d.Client.Disconnect(ctx)
}

func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		log.Fatal("Database not initialized. Call InitializeMongoDB() first.")
	}
	return DB.Collection(collectionName)
}

func IsConnected() bool {
    if Client == nil || DB == nil {
        return false
    }
    return true
}
