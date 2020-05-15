package linkstatusmongo

import (
	"encoding/json"

	hub "github.com/stottle-uk/my-first-go-app/internal/features/hub"
	domains "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/domains"
	storage "github.com/stottle-uk/my-first-go-app/internal/services/storage"
)

// LinkStatusRepo : LinkStatusRepo
type LinkStatusRepo struct {
	links *storage.Collection
	hub   *hub.Hub
}

// Options : Options
type Options struct {
	Links *storage.Collection
	Hub   *hub.Hub
}

// New : New
func New(options Options) *LinkStatusRepo {
	return &LinkStatusRepo{
		links: options.Links,
		hub:   options.Hub,
	}
}

// Insert : Insert
func (repo *LinkStatusRepo) Insert(doc domains.LinkStatus) (string, error) {
	insertID, err := repo.links.InsertOne(doc)
	if err != nil {
		return "", err
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}

	repo.hub.SendByTaskID <- &hub.SendByTask{
		TaskID:  doc.TaskID,
		Message: docBytes,
	}

	return insertID, err
}
