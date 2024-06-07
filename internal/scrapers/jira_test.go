package scrapers

import (
	"fmt"
	"testing"
	// "github.com/Sam-Pewton/prodex/internal/scrapers/types"
	"github.com/Sam-Pewton/prodex/internal/logging"
	"github.com/jarcoal/httpmock"
)

var logLevel, _ = logging.GetLogLevel("debug")

func TestNewScraper(t *testing.T) {
	scraper := NewJiraScraper(
		map[string]any{
			"AtlassianToken":  "abc123",
			"AtlassianUser":   "test@testemail.com",
			"AtlassianDomain": "https://test.atlassian.net",
			"PaginationSize":  50,
		},
	)

	// make sure that all of the members of the scraper are correctly set
	if scraper.headers["Accept"] != "application/json" {
		t.Fatalf("the headers for the scraper are incorrect")
	}
	if scraper.config["AtlassianToken"] != "abc123" {
		t.Fatalf("the atlassian token for the scraper is incorrect")
	}
	if scraper.config["AtlassianUser"] != "test@testemail.com" {
		t.Fatalf("the atlassian user for the scraper is incorrect")
	}
	if scraper.config["AtlassianDomain"] != "https://test.atlassian.net" {
		t.Fatalf("the atlassian domain for the scraper is incorrect")
	}
	if scraper.config["PaginationSize"] != 50 {
		t.Fatalf("the pagination size for the scraper is incorrect")
	}
}

func TestRetrieveAccountData(t *testing.T) {
	logging.StartLogger(logLevel)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://test.atlassian.net/rest/api/3/user/search?query=test@testemail.com",
		httpmock.NewStringResponder(
			200,
			`[{"accountId": "12345", "displayName": "test"}]`,
		),
	)

	scraper := NewJiraScraper(
		map[string]any{
			"AtlassianToken":  "abc123",
			"AtlassianUser":   "test@testemail.com",
			"AtlassianDomain": "https://test.atlassian.net",
			"PaginationSize":  50,
		},
	)

	user, err := scraper.retrieveAccountData()

	if err != nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}

	if user.AccountId != "12345" || user.DisplayName != "test" {
		t.Fatalf("the response payload was not properly unmarshalled")
	}
}

func TestInvalidStatusCode(t *testing.T) {
	logging.StartLogger(logLevel)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://test.atlassian.net/rest/api/3/user/search?query=test@testemail.com",
		httpmock.NewStringResponder(400, ``),
	)

	scraper := NewJiraScraper(
		map[string]any{
			"AtlassianToken":  "abc123",
			"AtlassianUser":   "test@testemail.com",
			"AtlassianDomain": "https://test.atlassian.net",
			"PaginationSize":  50,
		},
	)

	_, err := scraper.retrieveAccountData()

	if err == nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}
}

func TestNoUsersInResponsePayloadErr(t *testing.T) {
	logging.StartLogger(logLevel)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://test.atlassian.net/rest/api/3/user/search?query=test@testemail.com",
		httpmock.NewStringResponder(
			200,
			`[]`,
		),
	)

	scraper := NewJiraScraper(
		map[string]any{
			"AtlassianToken":  "abc123",
			"AtlassianUser":   "test@testemail.com",
			"AtlassianDomain": "https://test.atlassian.net",
			"PaginationSize":  50,
		},
	)

	_, err := scraper.retrieveAccountData()

	if err == nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}
}

func TestTwoUsersInResponsePayloadErr(t *testing.T) {
	logging.StartLogger(logLevel)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://test.atlassian.net/rest/api/3/user/search?query=test@testemail.com",
		httpmock.NewStringResponder(
			200,
			`[{"accountId": "12345", "displayName": "test"}, {"accountId": "67890", "displayName": "test2"}]`,
		),
	)

	scraper := NewJiraScraper(
		map[string]any{
			"AtlassianToken":  "abc123",
			"AtlassianUser":   "test@testemail.com",
			"AtlassianDomain": "https://test.atlassian.net",
			"PaginationSize":  50,
		},
	)

	_, err := scraper.retrieveAccountData()

	if err == nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}
}

func TestProcessData(t *testing.T) {
	logging.StartLogger(logLevel)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://test.atlassian.net/rest/api/3/search?jql=test=12345&startAt=0&maxResults=50",
		httpmock.NewStringResponder(
			200,
			`{
                "issues": [
                    {
                        "id": "1",
                        "key": "A",
                        "self": "abc",
                        "fields": {
                            "statuscategorychangedate": "20240607",
                            "issuetype": {
                                "name": "task"
                            },
                            "parent": {
                                "id": "2",
                                "key": "A",
                                "fields": {
                                    "summary": "test",
                                    "status": {
                                        "name": "test"
                                    },
                                    "priority": {
                                        "name": "high"
                                    },
                                    "issuetype": {
                                        "id": "3",
                                        "description": "aoeu",
                                        "name": "aeu"
                                    }
                                }
                            },
                            "project": {
                                "id": "4",
                                "key": "A",
                                "name": "test"
                            },
                            "created": "20240607",
                            "priority": {
                                "name": "high"
                            },
                            "issuelinks": [],
                            "assignee": {
                                "accountId": "12345",
                                "displayName": "test"
                            },
                            "updated": "20240607",
                            "status": {
                                "name": "test"
                            },
                            "summary": "test",
                            "creator": {
                                "accountId": "12345",
                                "displayName": "test"
                            },
                            "reporter": {
                                "accountId": "12345",
                                "displayName": "test"
                            },
                            "resolution": {
                                "name": "done"
                            },
                            "resolutiondate": "20240607T15:04:05+0100",
                            "description": {
                                "type": "test",
                                "content": []
                            }
                        }
                    }
                ]
            }`,
		),
	)

	scraper := NewJiraScraper(
		map[string]any{
			"AtlassianToken":  "abc123",
			"AtlassianUser":   "test@testemail.com",
			"AtlassianDomain": "https://test.atlassian.net",
			"PaginationSize":  50,
		},
	)

	testChan := make(chan DBExecution, 1)

	err := scraper.processData("test", "12345", testChan)

	if err != nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}

	data := <-testChan

	if data == nil {
		t.Fatalf(fmt.Sprintf("%s", data))
	}

	jiraData, ok := data.(JiraDBExecution)
	if !ok {
		t.Fatalf("data received was not a DB execution")
	}

	if jiraData.Table != "jira" {
		t.Fatalf(fmt.Sprintf("expected `jira` table, got `%s`", jiraData.Table))
	}
	if jiraData.Values[0] != "A" {
		t.Fatalf(fmt.Sprintf("%s", jiraData.Values))
	}
}
