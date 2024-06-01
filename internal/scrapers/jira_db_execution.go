// Jira DB execution related structs and methods.
package scrapers

import (
	"database/sql"
	"fmt"
	"github.com/Sam-Pewton/prodex/internal/logging"
)

// A DB execution struct for Jira.
type JiraDBExecution struct {
	PrimaryKeyIndex rune
	Table           string
	Values          [22]string
}

// Simply (and lazily) replace a record if it exists, or insert a new record.
// This will ALWAYS replace.
func (j JiraDBExecution) InsertOrUpdate(db *sql.DB) {
	logging.Debug(fmt.Sprintf("Inserting or updating record for %s", j.Values[j.PrimaryKeyIndex]))
	db.Exec(
		fmt.Sprintf(
			"REPLACE INTO %s VALUES (%s);",
			j.Table,
			valueArrayToString(j.Values[:]),
		),
	)
	logging.Info(fmt.Sprintf("Record inserted or updated for %s", j.Values[j.PrimaryKeyIndex]))
}
