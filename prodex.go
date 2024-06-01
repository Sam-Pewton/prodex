package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"github.com/Sam-Pewton/prodex/internal/config"
	"github.com/Sam-Pewton/prodex/internal/logging"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

//go:embed dev/schema.sql
var schema string

func setUpConfig() (*sql.DB, error) {
	err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	logLevel, err := logging.GetLogLevel(config.ProdexConf.LogLevel)
	logging.StartLogger(logLevel)
	if err != nil {
		// Runs after starting the logger, as need the logger in order to log..
		logging.Error(fmt.Sprintf("%s", err))
	}

	logging.Debug("Setting up DB")
	// Ensure that the required directory exists
	err = os.MkdirAll(
		fmt.Sprintf("%s%sdb/", os.Getenv("HOME"), config.ProdexConf.InstallationPath), 0755,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db, err := sql.Open(
		"sqlite3",
		fmt.Sprintf("%s/%s/db/prodex.db", os.Getenv("HOME"), config.ProdexConf.InstallationPath),
	)
	if err != nil {
		logging.Error("Could not open sqlite3 database")
		fmt.Println(err)
		return nil, err
	}
	logging.Debug("Connected to DB")

	if _, err := db.Exec(schema); err != nil {
		logging.Error(fmt.Sprint(err))
		db.Close()
		//fmt.Println(err)
		return nil, err
	}
	logging.Debug("DB created")
	return db, nil
}

func main() {
	var modeFlag = flag.String("mode", "scraper", "Mode of operation [scraper, ui]")
	flag.Parse()

	db, err := setUpConfig()
	if err != nil {
		return
	}
	defer db.Close()

	switch *modeFlag {
	case "scraper":
		runScraper(db)
	case "ui":
		runUI(db)
	default:
		logging.Error(fmt.Sprintf("unknown mode of operation: `%s`", *modeFlag))
	}
}
