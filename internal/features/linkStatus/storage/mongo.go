package linkstatusmongo

import (
	"encoding/json"
	"strconv"

	domains "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/domains"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
	storage "github.com/stottle-uk/my-first-go-app/internal/services/storage"
	"gopkg.in/mgo.v2/bson"
)

// LinkStatusRepo : LinkStatusRepo
type LinkStatusRepo struct {
	links *storage.Collection
	tasks *storage.Collection
	hub   *wshub.Hub
}

// GetTaskQuery : GetTaskQuery
type GetTaskQuery struct {
	TaskID int
}

// New : New
func New(db *storage.Database, hub *wshub.Hub) *LinkStatusRepo {
	return &LinkStatusRepo{
		links: db.Collection("links"),
		tasks: db.Collection("tasks"),
		hub:   hub,
	}
}

// Insert : Insert
func (repo *LinkStatusRepo) Insert(doc domains.LinkStatus) (string, error) {
	insertID, err := repo.links.InsertOne(doc)
	if err != nil {
		return "", err
	}

	task, err := repo.FindTask(&GetTaskQuery{
		TaskID: doc.TaskID,
	})
	if err != nil {
		return "", err
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}

	repo.hub.SendRestricted <- wshub.Message{
		ClientIds: []string{strconv.Itoa(task.UserID)},
		Data:      docBytes,
	}

	return insertID, err
}

// FindTask : FindTask
func (repo *LinkStatusRepo) FindTask(query *GetTaskQuery) (*domains.Task, error) {
	task := &domains.Task{}
	err := repo.tasks.FindOne(bson.M{
		"taskId": query.TaskID,
	}, &task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
