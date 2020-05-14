package linkstatusapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	domains "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/domains"
	store "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/storage"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

// API : API
type API struct {
	hub   *wshub.Hub
	store *store.LinkStatusRepo
}

// Options : Options
type Options struct {
	Hub   *wshub.Hub
	Store *store.LinkStatusRepo
}

// NewAPI : NewAPI
func NewAPI(options Options) (*API, error) {
	s := &API{
		hub:   options.Hub,
		store: options.Store,
	}
	return s, nil
}

// AddLink : AddLink
func (s *API) AddLink(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	linkStatus := domains.LinkStatus{}
	if err := json.Unmarshal(body, &linkStatus); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	linkStatusAdmin := domains.LinkStatusRequestAdmin{
		Data: domains.LinkStatusAdmin{
			URL:         linkStatus.URL,
			PageFoundOn: linkStatus.PageFoundOn,
			TaskID:      linkStatus.TaskID,
			ProductID:   linkStatus.ProductID,
		},
	}

	requestData, err := json.Marshal(linkStatusAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	z := bytes.NewBuffer(requestData)
	proxyReq, err := http.NewRequest(r.Method, "http://localhost:3333/api/links", z)

	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	res, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write body back to response
	body2, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ins, err := s.store.Insert(linkStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Insert ID: %v", ins)

	_, err = w.Write(body2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write header
	for header, values := range res.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.Header().Add("Content-Type", "application/json")

	go func() {
		s.hub.SendRestricted <- wshub.Message{ClientIds: []string{"cliendId123"}, Data: body}
	}()
}
