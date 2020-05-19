package linkstatusmongo

import (
	"encoding/json"

	hub "github.com/stottle-uk/my-first-go-app/internal/features/hub"
	domains "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/domains"
)

// LinkStatusHub : LinkStatusHub
type LinkStatusHub struct {
	hub     *hub.Hub
	checker domains.LinkChecker
}

// Options : Options
type Options struct {
	Hub     *hub.Hub
	Checker domains.LinkChecker
}

// New : New
func New(options Options) *LinkStatusHub {
	return &LinkStatusHub{
		hub:     options.Hub,
		checker: options.Checker,
	}
}

// AddLink : AddLink
func (repo *LinkStatusHub) AddLink(doc domains.LinkStatus) (string, error) {
	insertID, err := repo.checker.AddLink(doc)
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
