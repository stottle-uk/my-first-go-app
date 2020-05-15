package storagemongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"gopkg.in/mgo.v2/bson"
)

// Database : Database
type Database struct {
	db *mongo.Database
}

// Collection : Collection
type Collection struct {
	col            *mongo.Collection
	defaultTimeout time.Duration
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
		db: client.Database(cs.Database),
	}
}

// Collection : Collection
func (s *Database) Collection(name string) *Collection {
	return &Collection{
		col:            s.db.Collection(name),
		defaultTimeout: 5 * time.Second,
	}
}

// InsertOne : InsertOne
func (c *Collection) InsertOne(doc interface{}) (string, error) {
	result, err := c.col.InsertOne(context.Background(), doc)
	if err != nil {
		return "", err
	}

	itemID := result.InsertedID.(primitive.ObjectID).Hex()

	return itemID, nil
}

// FindOne : FindOne
func (c *Collection) FindOne(filter, out interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.defaultTimeout)
	defer cancel()

	return c.col.FindOne(ctx, filter).Decode(out)
}

// GetByID : GetByID
func (c *Collection) GetByID(itemID string, out interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return err
	}

	return c.FindOne(bson.M{"_id": objectID}, out)
}
