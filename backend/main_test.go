package main

import (
	"github.com/joho/godotenv"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

var a App

func init() {
	err := godotenv.Load()
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Problem with loading .env file")
	}
}

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("POSTGRES_NONROOT_USER"),
		os.Getenv("POSTGRES_NONROOT_PASSWORD"),
		os.Getenv("POSTGRES_NONROOT_DB"))

	ensureTableExists()
	code := m.Run()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal()
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS threads (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			link TEXT NOT NULL,
			posts INTEGER NOT NULL,
			votes INTEGER NOT NULL,
			views INTEGER NOT NULL,
			date_posted TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			seen BOOLEAN NOT NULL DEFAULT FALSE
		);`

func TestTable(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/deals", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
