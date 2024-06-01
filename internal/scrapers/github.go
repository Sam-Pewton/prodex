package scrapers

import (
	"errors"
	"fmt"
	"github.com/Sam-Pewton/prodex/internal/logging"
	"io"
	"net/http"
	"time"
)

type gitRepo struct {
	nodeId   string
	name     string
	fullName string
}

type gitScraper struct {
	url     string
	client  http.Client
	headers map[string]string
}

func NewGitScraper(url string, headers map[string]string) *gitScraper {
	return &gitScraper{
		url,
		http.Client{Timeout: time.Duration(1) * time.Second},
		headers,
	}
}

func (s *gitScraper) BuildGetRequest(url_suffix string) (*http.Request, error) {
	req, err := http.NewRequest("GET", s.url+url_suffix, nil)
	if err != nil {
		logging.Error("error %s", err)
	}

	for k, v := range s.headers {
		req.Header.Add(k, v)
	}

	return req, nil
}

func (s *gitScraper) ExecuteRequest(request *http.Request) (*http.Response, error) {
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

func (s *gitScraper) Scrape(c chan []byte) error {
	// 1. Check that we can access the URL via the /meta endpoint
	req, err := s.BuildGetRequest("")
	if err != nil {
		logging.Error("error %s", err)
		return err
	}

	res, err := s.ExecuteRequest(req)
	if err != nil {
		logging.Error("error %s", err)
		return err
	}

	// 2. Find all of the repository names
	req, err = s.BuildGetRequest("users/Sam-Pewton/repos")
	if err != nil {
		logging.Error("error %s", err)
		return err
	}

	res, err = s.ExecuteRequest(req)
	if err != nil {
		logging.Error("error %s", err)
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	c <- body
	return nil
}
