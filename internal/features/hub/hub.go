package hub

import (
	"net/http"
	"strconv"

	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
	storage "github.com/stottle-uk/my-first-go-app/internal/services/storage"

	"gopkg.in/mgo.v2/bson"
)

// Task : Task
type Task struct {
	TaskID int `json:"taskId"`
	UserID int `json:"userId"`
}

// SendByTask : SendByTask
type SendByTask struct {
	TaskID  int
	Message []byte
}

// Hub : Hub
type Hub struct {
	base         *wshub.Hub
	tasks        *storage.Collection
	SendByTaskID chan *SendByTask
	Broadcast    chan []byte
	Received     chan []byte
}

// GetTaskQuery : GetTaskQuery
type GetTaskQuery struct {
	TaskID int
}

// NewHub : NewHub
func NewHub(col *storage.Collection) (*Hub, func(w http.ResponseWriter, r *http.Request)) {
	h, handler := wshub.CreateHub()

	hub := &Hub{
		base:         h,
		tasks:        col,
		SendByTaskID: make(chan *SendByTask),
		Broadcast:    make(chan []byte),
		Received:     make(chan []byte),
	}

	go hub.run()

	return hub, handler
}

func (h *Hub) run() {
	for {
		select {
		case doc := <-h.SendByTaskID:
			clientIds, err := h.getClientIdsByTask(doc.TaskID)
			if err != nil {
				break
			}
			h.base.SendRestricted <- &wshub.Message{
				ClientIds: clientIds,
				Data:      doc.Message,
			}

		case message := <-h.Broadcast:
			h.base.Broadcast <- message

		case message := <-h.base.Received:
			h.Received <- message
		}
	}
}

func (h *Hub) getClientIdsByTask(taskID int) ([]string, error) {
	task, err := h.findTask(&GetTaskQuery{
		TaskID: taskID,
	})
	if err != nil {
		return nil, err
	}

	return []string{strconv.Itoa(task.UserID)}, nil
}

func (h *Hub) findTask(query *GetTaskQuery) (*Task, error) {
	task := &Task{}
	err := h.tasks.FindOne(bson.M{
		"taskId": query.TaskID,
	}, &task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
