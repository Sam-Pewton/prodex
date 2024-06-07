package main

import (
	"database/sql"
	"fmt"
	"github.com/Sam-Pewton/prodex/internal/config"
	"github.com/Sam-Pewton/prodex/internal/logging"
	"github.com/Sam-Pewton/prodex/internal/scrapers"
)

func runScraper(db *sql.DB) {
	totalScrapers := len(config.ProdexConf.Scrapers)

	// There are no scrapers to run.
	if totalScrapers == 0 {
		logging.Debug("There are no scrapers to run")
		return
	}

	// Data channels
	tasks := make(chan scrapers.DBExecution)
	orchestrator := make(chan bool, totalScrapers)
	defer close(tasks)
	defer close(orchestrator)

	// Start all scrapers
	for k, v := range config.ProdexConf.Scrapers {
		scrapers.MapToScraper(tasks, orchestrator, k, v)
	}

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

		if scrapersFinished && noop >= config.ProdexConf.MaxNoops {
			break
		}

		if noop >= config.ProdexConf.MaxNoops {
			logging.Error(fmt.Sprintf("scrapers stopped without finishing. expected %d, but only %d finished succesfully", totalScrapers, totalDone))
		}
	}
}
