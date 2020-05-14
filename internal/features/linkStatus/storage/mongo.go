package linkstatusmongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domains "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/domains"
)

// LinkStatusRepo : LinkStatusRepo
type LinkStatusRepo struct {
	collection     *mongo.Collection
	defaultTimeout time.Duration
}

// New : New
func New(db *mongo.Database) *LinkStatusRepo {
	return &LinkStatusRepo{
		collection:     db.Collection("links"),
		defaultTimeout: 5 * time.Second,
	}
}

// Insert : Insert
func (repo *LinkStatusRepo) Insert(doc domains.LinkStatus) (string, error) {
	result, err := repo.collection.InsertOne(context.Background(), doc)
	if err != nil {
		return "", err
	}

	itemID := result.InsertedID.(primitive.ObjectID).Hex()

	return itemID, nil
}
