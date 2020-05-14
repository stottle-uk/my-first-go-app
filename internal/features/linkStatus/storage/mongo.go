package linkstatusmongo

import (
	domains "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/domains"
	storage "github.com/stottle-uk/my-first-go-app/internal/services/storage"
)

// LinkStatusRepo : LinkStatusRepo
type LinkStatusRepo struct {
	col *storage.Collection
}

// New : New
func New(db *storage.Database) *LinkStatusRepo {
	return &LinkStatusRepo{
		col: db.Collection("links"),
	}
}

// Insert : Insert
func (repo *LinkStatusRepo) Insert(doc domains.LinkStatus) (string, error) {
	return repo.col.InsertOne(doc)
}
