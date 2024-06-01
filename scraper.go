package main

import (
	"database/sql"
	"fmt"
	"github.com/Sam-Pewton/prodex/internal/logging"
	"github.com/Sam-Pewton/prodex/internal/scrapers"
	"os"
)

func runScraper(db *sql.DB) {

	// Data channels
	tasks := make(chan scrapers.DBExecution)
	orchestrator := make(chan bool, 10)
	defer close(tasks)
	defer close(orchestrator)

	// go scrapeGitHub(tasks, orchestrator)
	// go scrapeConfluence(tasks, orchestrator)
	go scrapeJira(tasks, orchestrator)

	totalDone := 0
	scrapersFinished := false
	noop := 0
	for {
		// Check if one of the scrapers is finished
		if len(orchestrator) > 0 {
			select {
			case <-orchestrator:
				totalDone++
				if totalDone >= totalScrapers {
					logging.Debug("All scrapers finished")
					scrapersFinished = true
				}
			}
		}

		// Try to extract some data from the queue
		select {
		case data := <-tasks:
			// handle the received data and reset the noop count
			noop = 0
			data.InsertOrUpdate(db)
		default:
			noop++
		}

		if scrapersFinished && noop >= maxNoops {
			break
		}
	}
}

func scrapeGitHub(c chan []byte, orchestrator chan<- bool) {
	headers := map[string]string{
		"Accept":               "application/vnd.github+json",
		"Authorization":        fmt.Sprintf("Bearer %s", os.Getenv("GH_TOKEN")),
		"X-GitHub-Api-Version": "2022-11-28",
	}

	s := scrapers.NewGitScraper(
		"https://api.github.com/",
		headers,
	)

	// Run the scraper
	logging.Debug("GitHub scraper starting")
	s.Scrape(c)

	// Signal to the main thread that this scraper has finished.
	orchestrator <- true
	logging.Debug("GitHub scraper finished")
}

func scrapeJira(c chan<- scrapers.DBExecution, orchestrator chan<- bool) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	s := scrapers.NewJiraScraper(
		os.Getenv("ATLASSIAN_DOMAIN"),
		headers,
	)

	// Run the scraper
	logging.Debug("Jira scraper starting")
	s.Scrape(c)

	// Signal to the main thread that this scraper has finished.
	orchestrator <- true
	logging.Debug("Jira scraper finished")
}

func scrapeConfluence(c chan<- []byte, orchestrator chan<- bool) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	s := scrapers.NewConfluenceScraper(
		os.Getenv("ATLASSIAN_DOMAIN"),
		headers,
	)

	// Run the scraper
	logging.Debug("Confluence scraper starting")
	s.Scrape(c)

	// Signal to the main thread that this scraper has finished.
	orchestrator <- true
	logging.Debug("Confluence scraper finished")
}
