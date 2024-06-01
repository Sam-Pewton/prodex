package main

import (
	"database/sql"
	"github.com/Sam-Pewton/prodex/internal/logging"
	"net/http"
)

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("<h1>Hello from Prodex!</h1>"))
}

func statsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("<h1>Hello from Prodex Stats!</h1>"))
}

func insertHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("<h1>Hello from Prodex insert!</h1>"))
}

func runUI(*sql.DB) {
	logging.Info("Hello from the UI")
	port := "8642"

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/stats", statsHandler)
	mux.HandleFunc("/insert", insertHandler)
	http.ListenAndServe(":"+port, mux)

}
