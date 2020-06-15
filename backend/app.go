package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
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
	RDB    *redis.Client
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
	connectRedis(a)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) handleDeals() http.HandlerFunc {
	log.Debug("Deals API endpoint registered")
	return func(w http.ResponseWriter, r *http.Request) {
		threads, err := getThreads(a)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "an error has occurred")
			log.Error(err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, threads)
	}
}

func (a *App) handleHealthCheck() http.HandlerFunc {
	log.Debug("Healthcheck API endpoint registered")
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a.DB.Ping(); err != nil {
			respondWithJSON(w, http.StatusInternalServerError, map[string]string{"message": "not ok"})
			log.Error(err.Error())
			return
		}
		var ctx = context.Background()
		if _, err := a.RDB.Ping(ctx).Result(); err != nil {
			respondWithJSON(w, http.StatusInternalServerError, map[string]string{"message": "not ok"})
			log.Error(err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "ok"})
	}
}

func (a *App) handleSubscribe() http.HandlerFunc {
	log.Debug("Subscribe API endpoint registered")
	return func(w http.ResponseWriter, r *http.Request) {
		var s subscriber
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&s); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		if err := s.createSubscriber(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusCreated, s)
	}
}

func (a *App) handleUnsubscribe() http.HandlerFunc {
	log.Debug("Unsubscribe API endpoint registered")
	return func(w http.ResponseWriter, r *http.Request) {
		var s subscriber
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&s); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		if err := s.deleteSubscriber(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, s)
	}
}

func (a *App) handleConfirm() http.HandlerFunc {
	log.Debug("Confirm API endpoint registered")
	return func(w http.ResponseWriter, r *http.Request) {
		var s subscriber
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&s); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		if err := s.updateSubscriber(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, s)
	}
}

func (a *App) initializeRoutes() {
	apiRoute := a.Router.PathPrefix("/api/v1").Subrouter()
	apiRoute.HandleFunc("/deals", a.handleDeals()).Methods("GET")
	apiRoute.HandleFunc("/healthcheck", a.handleHealthCheck()).Methods("GET")
	emailsRoute := apiRoute.PathPrefix("/emails").Subrouter()
	emailsRoute.HandleFunc("", a.handleSubscribe()).Methods("POST")
	emailsRoute.HandleFunc("", a.handleUnsubscribe()).Methods("DELETE")
	emailsRoute.HandleFunc("", a.handleConfirm()).Methods("PUT")
}
