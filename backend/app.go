package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func init() {
	formatter := &log.TextFormatter{FullTimestamp: true}
	log.SetFormatter(formatter)
	log.SetLevel(log.DebugLevel)
	err := godotenv.Load()
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Problem with loading .env file")
	}
}

func (a *App) Initialize(user, password, dbname string) {
	host, exists := os.LookupEnv("POSTGRES_HOST")
	if !exists {
		host = "postgres"
	}
	port := 5432
	pgURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	connectDB(a, pgURI)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) handleDeals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		threads, err := getThreads(a.DB)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, threads)
	}
}

func (a *App) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := a.DB.Ping()
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, map[string]string{"message": "not ok"})
			log.Error(err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "ok"})
	}
}

func (a *App) initializeRoutes() {
	apiRoute := a.Router.PathPrefix("/api/v1").Subrouter()
	apiRoute.HandleFunc("/deals", a.handleDeals()).Methods("GET")
	apiRoute.HandleFunc("/healthcheck", a.handleHealthCheck()).Methods("GET")
}
// todo add more logging
