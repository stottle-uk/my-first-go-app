package linkstatusapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	domains "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/domains"
	store "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/storage"
	redirect "github.com/stottle-uk/my-first-go-app/internal/services/redirect"
)

// API : API
type API struct {
	store    *store.LinkStatusRepo
	redirect *redirect.Redirect
}

// Options : Options
type Options struct {
	Store    *store.LinkStatusRepo
	Redirect *redirect.Redirect
}

// NewAPI : NewAPI
func NewAPI(options Options) *API {
	return &API{
		store:    options.Store,
		redirect: options.Redirect,
	}
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

	res, err := s.redirect.Do("/links", requestData, r)

	bodyAdmin, err := ioutil.ReadAll(res.Body)
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

	_, err = w.Write(bodyAdmin)
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
}
