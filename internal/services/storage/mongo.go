package storagemongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// NewDb : NewDb
func NewDb() (*mongo.Database, error) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cs, err := connstring.Parse("mongodb://localhost:27017/testing123")
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cs.String()))
	return client.Database(cs.Database), err
}
