package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Problem with loading .env file")
	}
	log.SetLevel(log.DebugLevel)
}

func (a *App) Initialize(user, password, dbname string) {
	host, exists := os.LookupEnv("POSTGRES_HOST")
	if !exists {
		host = "postgres"
	}
	port := 5432
	pgURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	log.Info("Attempting to connect to DB")
	for i := 1; i < 6; i++ {
		a.DB, err = sql.Open("postgres", pgURI)
		if err != nil {
			log.Error("Error with opening connection with DB")
			panic(err)
		}

		err = a.DB.Ping()
		if err == nil {
			break
		}
		retryWait := i * i
		log.Info(fmt.Sprintf("Connection attempt %d unsuccessful, retrying in %d seconds...", i, retryWait))
		time.Sleep(time.Duration(retryWait) * time.Second)
	}
	if a.DB == nil {
		log.Fatal("Could not connect to DB")
	}
	log.Info("Successfully connected to DB")
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := getThreads(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, threads)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn()
	}

}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/v1/deals", a.getThreads).Methods("GET")
}
