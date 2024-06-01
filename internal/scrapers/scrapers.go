package scrapers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Sam-Pewton/prodex/internal/config"
	"github.com/Sam-Pewton/prodex/internal/logging"
)

// Scraper interface. All scrapers should use this.
type Scraper interface {
	Scrape(chan<- []byte)
	buildGetRequest(string) (*http.Request, error)
	executeRequest(*http.Request) error
}

// Common type used for inserting to the DB.
// Main thread will run the associated method.
type DBExecution interface {
	InsertOrUpdate(*sql.DB)
}

// Transform a slice of strings into a value string.
// Used on insert/update/replace statements.
func valueArrayToString(values []string) string {
	v := ""
	for _, x := range values {
		if v == "" {
			v = "'" + x + "'"
		} else {
			v = v + ", '" + x + "'"
		}
	}
	return v
}

func MapToScraper(tasks chan<- DBExecution, orchestrator chan<- bool, scraperName string, configs config.ScraperConfigs) {
	for _, c := range configs {
		// scrapers.MapToScraper(tasks, orchestrator, k, c)
		switch scraperName {
		case "jira":
			go scrapeJira(tasks, orchestrator, c)
		case "confluence":
			go scrapeConfluence(tasks, orchestrator, c)
		case "github":
			go scrapeGitHub(tasks, orchestrator, c)
		default:
			logging.Error(fmt.Sprintf("unknown scraper provided: %s", scraperName))
		}
	}
}

func scrapeGitHub(c chan<- DBExecution, orchestrator chan<- bool, config map[string]any) {
	s := NewGitScraper(config)

	// Run the scraper
	logging.Debug("GitHub scraper starting")
	s.Scrape(c)

	// Signal to the main thread that this scraper has finished.
	orchestrator <- true
	logging.Debug("GitHub scraper finished")
}

func scrapeJira(c chan<- DBExecution, orchestrator chan<- bool, config map[string]any) {
	s := NewJiraScraper(config)

	// Run the scraper
	logging.Debug("Jira scraper starting")
	s.Scrape(c)

	// Signal to the main thread that this scraper has finished.
	orchestrator <- true
	logging.Debug("Jira scraper finished")
}

func scrapeConfluence(c chan<- DBExecution, orchestrator chan<- bool, config map[string]any) {
	s := NewConfluenceScraper(config)

	// Run the scraper
	logging.Debug("Confluence scraper starting")
	s.Scrape(c)

	// Signal to the main thread that this scraper has finished.
	orchestrator <- true
	logging.Debug("Confluence scraper finished")
}
