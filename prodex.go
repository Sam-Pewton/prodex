package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Sam-Pewton/prodex/internal/logging"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed dev/schema.sql
var schema string

// The total number of scrapers that are running in the app
const totalScrapers = 1

// The maximum amount of noops allowed after all scrapers complete
const maxNoops = 10000

func loadDotEnv() error {
	err := godotenv.Load(os.Getenv("HOME") + "/.config/prodex/.env")
	if err != nil {
		logging.Error(fmt.Sprintf("error loading .env file from config %s. Attempting load from cwd", os.Getenv("HOME")+"/.config/prodex/.env"))
		err = godotenv.Load()
		if err != nil {
			logging.Error("error loading .env file from cwd. Exiting")
			return err
		}
	}
	return nil
}

func setUpConfig() (*sql.DB, error) {
	logging.Debug("Setting up DB")

	// Ensure that the required directory exists
	err := os.MkdirAll(fmt.Sprintf("%s%sdb/", os.Getenv("HOME"), os.Getenv("INSTALL_PATH")), 0755)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s/db/prodex.db", os.Getenv("HOME"), os.Getenv("INSTALL_PATH")))
	if err != nil {
		logging.Error("Could not open sqlite3 database")
		fmt.Println(err)
		return nil, err
	}
	logging.Debug("DB opened")

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

	logging.StartLogger(slog.LevelInfo)
	err := loadDotEnv()
	if err != nil {
		return
	}

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
		logging.Error(fmt.Sprintf("unknown mode of operation: %s", *modeFlag))
	}
}
