package scrapers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Sam-Pewton/prodex/internal/logging"
)

type confluenceScraper struct {
	client  http.Client
	headers map[string]string
	config  map[string]any
}

func NewConfluenceScraper(config map[string]any) *confluenceScraper {
	headers := map[string]string{
		"Accept": "application/json",
	}
	return &confluenceScraper{
		http.Client{Timeout: time.Duration(1) * time.Second},
		headers,
		config,
	}
}

func (s *confluenceScraper) BuildGetRequest(url_suffix string) (*http.Request, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%s", s.config["atlassian_domain"], url_suffix),
		nil,
	)
	if err != nil {
		logging.Error("error %s", err)
	}

	for k, v := range s.headers {
		req.Header.Add(k, v)
	}

	req.SetBasicAuth(
		fmt.Sprintf("%s", s.config["atlassian_user"]),
		fmt.Sprintf("%s", s.config["atlassian_token"]),
	)

	return req, nil
}

func (s *confluenceScraper) ExecuteRequest(request *http.Request) (*http.Response, error) {
	resp, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	fmt.Println(resp.Status)
	if resp.Status != "200 OK" {
		return nil, errors.New("error: invalid status code received from meta endpoint. Cannot continue")
	}
	return resp, nil
}

func (s *confluenceScraper) Scrape(c chan<- DBExecution) error {
	// 1. Check that we can access the URL via the /meta endpoint
	req, err := s.BuildGetRequest("wiki/api/v2/spaces")
	if err != nil {
		logging.Error("error %s", err)
		return err
	}

	// 2. Use the same endpoint as in Jira to get the user ID
	//    I can't see a way to use the Confluence API to do this

	// 3. Get the space ID
	//    This will potentially be done by filtering on the key
	//    https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-space/#api-spaces-get-request

	// For each of the spaces in the array, I should spawn a goroutine to scrape individually

	res, err := s.ExecuteRequest(req)
	if err != nil {
		logging.Error("error %s", err)
		return err
	}

	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)

	//c <- body
	return nil
}
