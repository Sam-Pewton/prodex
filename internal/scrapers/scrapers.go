package scrapers

import (
	"database/sql"
	"net/http"
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
