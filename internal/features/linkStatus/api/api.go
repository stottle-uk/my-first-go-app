package linkstatusapi

import (
	"io/ioutil"
	"net/http"

	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

// API : API
type API struct {
	hub *wshub.Hub
}

// Options : Options
type Options struct {
	Hub *wshub.Hub
}

// NewAPI : NewAPI
func NewAPI(options Options) (*API, error) {
	s := &API{
		hub: options.Hub,
	}
	return s, nil
}

// AddLink : AddLink
func (s *API) AddLink(w http.ResponseWriter, r *http.Request) {
	proxyReq, err := http.NewRequest(r.Method, "http://localhost:3333/api/links", r.Body)

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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
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
