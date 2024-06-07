// Jira scraper creator and methods
package scrapers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sam-Pewton/prodex/internal/logging"
	"github.com/Sam-Pewton/prodex/internal/scrapers/types"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// A jira scraper.
type jiraScraper struct {
	client  http.Client
	headers map[string]string
	config  map[string]any
}

// Create a new jira scraper.
func NewJiraScraper(config map[string]any) *jiraScraper {
	headers := map[string]string{
		"Accept": "application/json",
	}

	return &jiraScraper{
		http.Client{Timeout: time.Duration(5) * time.Second},
		headers,
		config,
	}
}

// Build a basic GET request to execute with the client
func (s *jiraScraper) buildGetRequest(url_suffix string) (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", s.config["AtlassianDomain"], url_suffix), nil)
	if err != nil {
		logging.Error("error %s", err)
	}

	for k, v := range s.headers {
		req.Header.Add(k, v)
	}
	req.SetBasicAuth(fmt.Sprintf("%s", s.config["AtlassianUser"]), fmt.Sprintf("%s", s.config["AtlassianToken"]))

	return req, nil
}

// Execute an HTTP request using the scrapers client.
func (s *jiraScraper) executeRequest(request *http.Request) (*http.Response, error) {
	resp, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	if !strings.Contains(resp.Status, "200") {
		return nil, errors.New("error: invalid status code received from meta endpoint. Cannot continue")
	}
	return resp, nil
}

// Retrieve the account data for a Jira user.
func (s *jiraScraper) retrieveAccountData() (*types.JiraUser, error) {
	req, err := s.buildGetRequest(fmt.Sprintf("rest/api/3/user/search?query=%s", s.config["AtlassianUser"]))

	if err != nil {
		return nil, err
	}

	res, err := s.executeRequest(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var users []types.JiraUser
	err = json.Unmarshal([]byte(body), &users)
	if err != nil {
		return nil, err
	}

	switch len(users) {
	case 0:
		return nil, fmt.Errorf("could not find an associated user.")
	case 1:
		return &users[0], nil
	default:
		return nil, fmt.Errorf("too many users found")
	}
}

// Process jira data for a specific query.
// queryType should be in `creator, assignee, reporter`
func (s *jiraScraper) processData(queryType string, accID string, c chan<- DBExecution) error {
	currentLen, ok := s.config["PaginationSize"].(int)
	if !ok {
		logging.Error("pagination size needs to be an integer")
	}
	currentPos := 0

	for {
		// Build the request
		req, err := s.buildGetRequest(
			fmt.Sprintf("rest/api/3/search?jql=%s=%s&startAt=%d&maxResults=%d", queryType, accID, currentPos, s.config["PaginationSize"]),
		)
		if err != nil {
			logging.Error(fmt.Sprint(err))
			return err
		}

		// Execute the request
		res, err := s.executeRequest(req)
		if err != nil {
			logging.Error(fmt.Sprint(err))
			return err
		}

		// Read the HTTP response body
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			logging.Error(fmt.Sprint(err))
			return err
		}

		// Unmarshall the response
		var issues types.JiraIssueResponse
		err = json.Unmarshal([]byte(body), &issues)
		if err != nil {
			logging.Error(fmt.Sprint(err))
			return err
		}

		currentLen = len(issues.Issues)

		// Push the data into the channel
		for _, record := range issues.Issues {
			db_exec := s.buildJiraDBExecution(&record)
			select {
			case c <- db_exec:
			}
		}

		// We have exhausted all content pages, we are done.
		if currentLen < s.config["PaginationSize"].(int) {
			break
		}

		// Move a page up
		currentPos += currentLen
	}
	return nil
}

// Convert a JIRA timestamp into that expected by SQLite3
func convertJiraDate(date string) string {
	converted, err := time.Parse("2006-01-02T15:04:05+0100", date)
	if err != nil {
		logging.Debug(fmt.Sprintf("unparsable timestamp received `%s`. Potentially intentional", date))
		return "nil"
	}
	return converted.Format("2006-01-02 15:04:05")
}

// Build a JiraDBExecution for the main thread to execute
func (s *jiraScraper) buildJiraDBExecution(data *types.JiraIssue) JiraDBExecution {
	values := [22]string{
		data.Key,
		data.Self,
		data.Fields.IssueType.Name,
		data.Fields.Summary,
		strings.Join(data.Fields.Description.Content, ""),
		data.Fields.Status.Name,
		data.Fields.Priority.Name,
		data.Fields.Resolution.Name,
		data.Fields.Assignee.DisplayName,
		data.Fields.Creator.DisplayName,
		data.Fields.Reporter.DisplayName,
		data.Fields.Parent.Key,
		data.Fields.Parent.Fields.IssueType.Name,
		data.Fields.Parent.Fields.Summary,
		data.Fields.Parent.Fields.Status.Name,
		data.Fields.Parent.Fields.Priority.Name,
		data.Fields.Project.Key,
		data.Fields.Project.Name,
		convertJiraDate(data.Fields.Created),
		convertJiraDate(data.Fields.Updated),
		convertJiraDate(data.Fields.ResolutionDate),
		convertJiraDate(data.Fields.StatusCategoryChangeDate),
	}

	return JiraDBExecution{0, "jira", values}
}

// Run the jira scraper.
func (s *jiraScraper) Scrape(c chan<- DBExecution) error {
	// 1. Get the account data
	user, err := s.retrieveAccountData()
	if err != nil {
		logging.Error(fmt.Sprint(err))
		return err
	}

	var wg sync.WaitGroup

	// 2. Run for all queries
	for _, t := range []string{"assignee", "reporter", "creator"} {
		wg.Add(1)
		go func(t string) {
			logging.Debug(fmt.Sprintf("Starting jira: %s\n", t))
			defer wg.Done()
			s.processData(t, user.AccountId, c)
			logging.Debug(fmt.Sprintf("Stopping jira: %s\n", t))
		}(t)
	}
	wg.Wait()
	return nil
}
