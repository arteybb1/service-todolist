package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database
var MongoClient *mongo.Client

func MongoConfig() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(AppConfig.MongoURI))
	if err != nil {
		log.Fatal("MongoDB connect error:", err)
	}

	// ตรวจสอบการเชื่อมต่อด้วย Ping
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping error:", err)
	}

	MongoClient = client
	MongoDB = client.Database("todoplus")
}
