package linkstatusapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	domains "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/domains"
	redirect "github.com/stottle-uk/my-first-go-app/internal/services/redirect"
)

// API : API
type API struct {
	redirect *redirect.Redirect
	checker  domains.LinkChecker
}

// Options : Options
type Options struct {
	Redirect *redirect.Redirect
	Checker  domains.LinkChecker
}

// NewAPI : NewAPI
func NewAPI(options Options) *API {
	return &API{
		checker:  options.Checker,
		redirect: options.Redirect,
	}
}

// AddLink : AddLink
func (s *API) AddLink(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	linkStatus := domains.LinkStatus{}
	if err := json.Unmarshal(body, &linkStatus); err != nil {
		s.handleError(w, err)
		return
	}

	res, err := s.addLinkAdmin(linkStatus, r)
	if err != nil {
		s.handleError(w, err)
		return
	}

	bodyAdmin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		s.handleError(w, err)
		return
	}

	ins, err := s.checker.AddLink(linkStatus)
	if err != nil {
		s.handleError(w, err)
		return
	}

	w.Header().Add("x-inserted-id", ins)
	for header, values := range res.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(res.StatusCode)

	_, err = w.Write(bodyAdmin)
	if err != nil {
		s.handleError(w, err)
		return
	}
}

func (s *API) addLinkAdmin(linkStatus domains.LinkStatus, r *http.Request) (*http.Response, error) {
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
		return nil, err
	}

	return s.redirect.Do("/links", requestData, r)
}

func (s *API) handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
