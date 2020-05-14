package storagemongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// Database : Database
type Database struct {
	Db *mongo.Database
}

// Collection : Collection
type Collection struct {
	Col *mongo.Collection
}

// NewDb : NewDb
func NewDb() *Database {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cs, err := connstring.Parse("mongodb://localhost:27017/testing123")
	if err != nil {
		log.Printf("Connection String Error: %v", err)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cs.String()))
	if err != nil {
		log.Printf("Mongo Connect Error: %v", err)
	}

	return &Database{
		Db: client.Database(cs.Database),
	}
}

// Collection : Collection
func (s *Database) Collection(name string) *Collection {
	return &Collection{Col: s.Db.Collection(name)}
}

// InsertOne : InsertOne
func (c *Collection) InsertOne(doc interface{}) (string, error) {
	result, err := c.Col.InsertOne(context.Background(), doc)
	if err != nil {
		return "", err
	}

	itemID := result.InsertedID.(primitive.ObjectID).Hex()

	return itemID, nil
}
