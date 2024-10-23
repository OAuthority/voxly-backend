package database

import (
    "context"
    "os"
    "sync"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
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
        username := os.Getenv("MONGO_USERNAME")
        password := os.Getenv("MONGO_PASSWORD")
       
        if mongoURI == "" || username == "" || password == "" {
            log.Fatal("MongoDB URI, Username, or Password is not set in environment variables")
        }

        clientOptions := options.Client().ApplyURI(mongoURI).SetAuth(options.Credential{
            Username: username,
            Password: password,
        })

        client, err = mongo.Connect(context.Background(), clientOptions)
        if err != nil {
            log.Fatal("Failed to connect to MongoDB:", err)
        }
    })

    return client, err
}

// Get access to specific collection (table) in the database
func GetCollection(collectionName string) *mongo.Collection {
    // lets get a client first, obviously >.<
    if client == nil {
        var err error
        client, err = GetClient()
        if err != nil {
            log.Fatal("Failed to initialize MongoDB client:", err)
        }
    }

    dbName := os.Getenv("DB_NAME")
    if dbName == "" {
        log.Fatal("Database name is not set in environment variables")
    }
    return client.Database(dbName).Collection(collectionName)
}