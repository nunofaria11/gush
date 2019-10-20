package storage

import (
	"context"
	"fmt"
	"gush/models"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

const (
	dbName                  = "gush-db"
	collName                = "url-list"
	defaultConnectionString = "mongodb://localhost:27017"
)

var collection *mongo.Collection

func getMongoDBConnectionString() string {
	var connectionString string

	url, ok := os.LookupEnv("MONGO_URL")
	if ok {
		user, _ := os.LookupEnv("MONGO_USER")
		password, _ := os.LookupEnv("MONGO_PASSWORD")
		port, _ := os.LookupEnv("MONGO_PORT")
		connectionString = fmt.Sprintf("mongodb://%v:%v@%v:%v", user, password, url, port)
	} else {
		connectionString = defaultConnectionString
	}

	log.Print("Mongo DB connection:", connectionString)
	return connectionString
}

func init() {

	connectionString := getMongoDBConnectionString()

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database(dbName).Collection(collName)
	log.Print("URL storage initialized")
}

// StoreURLInfo stores a URL info in database
func StoreURLInfo(urlInfo *models.URLInfo) error {
	result, err := collection.InsertOne(context.Background(), &urlInfo)

	if err != nil {
		return fmt.Errorf("Couldn't store URL info: %v", err)
	}

	log.Print("Inserted a URL in database ", result.InsertedID)
	return nil
}

// FetchURLInfo fetches a URL info from database
func FetchURLInfo(hash string) (*models.URLInfo, error) {
	var urlInfo models.URLInfo

	err := collection.FindOne(context.Background(), bson.M{"hash": hash}).Decode(&urlInfo)
	if err != nil {
		return nil, fmt.Errorf("Couldn't fetch URL info: %v", err)
	}

	log.Printf("Fetched URL from database: %v", urlInfo)
	return &urlInfo, nil
}
