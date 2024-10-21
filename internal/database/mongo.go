package database

import (
    "context"
    "os"
    "sync"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    client *mongo.Client
    once   sync.Once
)

// Get a connection to the client (database)
// so that we may do funs stuff such as query, obviously!
func GetClient() (*mongo.Client, error) {
    var err error
    
    once.Do(func() {
        mongoURI := os.Getenv("MONGO_URI")
        clientOptions := options.Client().ApplyURI(mongoURI)
        client, err = mongo.Connect(context.Background(), clientOptions)
    })
    
    return client, err
}

// Get access to specific collection (table) in the database
func GetCollection(collectionName string) *mongo.Collection {
    dbName := os.Getenv("DB_NAME")
    return client.Database(dbName).Collection(collectionName)
}