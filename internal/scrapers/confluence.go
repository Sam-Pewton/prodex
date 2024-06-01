package scrapers

import (
	"errors"
	"fmt"
	"github.com/Sam-Pewton/prodex/internal/logging"
	"io"
	"net/http"
	"os"
	"time"
)

type confluenceScraper struct {
	url     string
	client  http.Client
	headers map[string]string
}

func NewConfluenceScraper(url string, headers map[string]string) *confluenceScraper {
	return &confluenceScraper{
		url,
		http.Client{Timeout: time.Duration(1) * time.Second},
		headers,
	}
}

func (s *confluenceScraper) BuildGetRequest(url_suffix string) (*http.Request, error) {
	req, err := http.NewRequest("GET", s.url+url_suffix, nil)
	if err != nil {
		logging.Error("error %s", err)
	}

	for k, v := range s.headers {
		req.Header.Add(k, v)
	}

	req.SetBasicAuth(os.Getenv("ATLASSIAN_USER"), os.Getenv("ATLASSIAN_TOKEN"))

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

func (s *confluenceScraper) Scrape(c chan<- []byte) error {
	// 1. Check that we can access the URL via the /meta endpoint
	req, err := s.BuildGetRequest("wiki/api/v2/spaces")
	if err != nil {
		logging.Error("error %s", err)
		return err
	}

	res, err := s.ExecuteRequest(req)
	if err != nil {
		logging.Error("error %s", err)
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	c <- body
	return nil
}
