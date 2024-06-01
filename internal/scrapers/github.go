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
	client  http.Client
	headers map[string]string
	config  map[string]any
}

func NewGitScraper(config map[string]any) *gitScraper {
	headers := map[string]string{
		"Accept":               "application/vnd.github+json",
		"Authorization":        fmt.Sprintf("Bearer %s", config["github_token"]),
		"X-GitHub-Api-Version": "2022-11-28",
	}
	return &gitScraper{
		http.Client{Timeout: time.Duration(1) * time.Second},
		headers,
		config,
	}
}

func (s *gitScraper) BuildGetRequest(url_suffix string) (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", s.config["github_url"], url_suffix), nil)
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

func (s *gitScraper) Scrape(c chan<- DBExecution) error {
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
	_, err = io.ReadAll(res.Body)

	// c <- body
	return nil
}
