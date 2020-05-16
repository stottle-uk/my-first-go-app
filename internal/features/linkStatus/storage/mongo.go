package linkstatusmongo

import (
	domains "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/domains"
	storage "github.com/stottle-uk/my-first-go-app/internal/services/storage"
)

// LinkStatusRepo : LinkStatusRepo
type LinkStatusRepo struct {
	links *storage.Collection
}

// Options : Options
type Options struct {
	Links *storage.Collection
}

// New : New
func New(options Options) *LinkStatusRepo {
	return &LinkStatusRepo{
		links: options.Links,
	}
}

// AddLink : AddLink
func (repo *LinkStatusRepo) AddLink(doc domains.LinkStatus) (string, error) {
	insertID, err := repo.links.InsertOne(doc)
	if err != nil {
		return "", err
	}

	return insertID, err
}
