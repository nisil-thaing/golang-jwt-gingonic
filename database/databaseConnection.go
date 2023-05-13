package database

import (
  "fmt"
  "log"
  "time"
  "os"
  "context"

  "github.com/joho/godotenv"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func GetDBInstance() *mongo.Client {
  if err := godotenv.Load(".env"); err != nil {
    log.Fatal("Loading .env file failed!") 
  }

  mongoUri := os.Getenv("MONGODB_URI")

  clientOption := options.Client().ApplyURI(mongoUri)
  client, err := mongo.NewClient(clientOption)

  if err != nil {
    log.Fatal(err)
  }

  ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
  defer cancel()

  if err := client.Connect(ctx); err != nil {
    log.Fatal(err)
  }

  fmt.Println("MongoDB's connection has been established recently!")

  return client
}

var Client *mongo.Client = GetDBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
  collection := client.Database("cluster0").Collection(collectionName)

  return collection
}
