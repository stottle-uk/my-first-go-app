package linkstatusmongo

import (
	"encoding/json"

	hub "github.com/stottle-uk/my-first-go-app/internal/features/hub"
	domains "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/domains"
	store "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/storage"
)

// LinkStatusHub : LinkStatusHub
type LinkStatusHub struct {
	hub   *hub.Hub
	store *store.LinkStatusRepo
}

// Options : Options
type Options struct {
	Hub   *hub.Hub
	Store *store.LinkStatusRepo
}

// New : New
func New(options Options) *LinkStatusHub {
	return &LinkStatusHub{
		store: options.Store,
		hub:   options.Hub,
	}
}

// AddLink : AddLink
func (repo *LinkStatusHub) AddLink(doc domains.LinkStatus) (string, error) {
	insertID, err := repo.store.AddLink(doc)
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
