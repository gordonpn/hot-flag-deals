package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func (a *App) connectDB() error {
	host, exists := os.LookupEnv("POSTGRES_HOST")
	if !exists {
		host = "postgres"
	}
	user := os.Getenv("POSTGRES_NONROOT_USER")
	password := os.Getenv("POSTGRES_NONROOT_PASSWORD")
	dbname := os.Getenv("POSTGRES_NONROOT_DB")
	port := 5432
	pgURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	log.Info("Attempting to connect to database")
	var err error
	for i := 1; i < 6; i++ {
		if a.DB, err = sql.Open("postgres", pgURI); err != nil {
			log.Error("Error with opening connection with DB")
			return err
		}

		if err = a.DB.Ping(); err == nil {
			break
		}
		retryWait := i * i
		log.Infof("Connection attempt %d unsuccessful, retrying in %d seconds...", i, retryWait)
		time.Sleep(time.Duration(retryWait) * time.Second)
	}
	if a.DB == nil {
		log.Fatal("Could not connect to database")
	}
	log.Info("Successfully connected to database")
	return nil
}

func (a *App) getThreads() error {
	log.Info("Fetching data from database")
	sqlQuery := `SELECT * FROM threads WHERE date_posted > CURRENT_TIMESTAMP - INTERVAL '30 day';`
	rows, err := a.DB.Query(sqlQuery)
	if err != nil {
		return err
	}

	for rows.Next() {
		row := thread{}
		err := rows.Scan(
			&row.ID,
			&row.Title,
			&row.Link,
			&row.Posts,
			&row.Votes,
			&row.Views,
			&row.DatePosted,
			&row.Seen,
			&row.Notified,
		)
		if err != nil {
			return err
		}
		a.threads = append(a.threads, row)
	}
	return nil
}

func (a *App) markNotified(threads []thread) error {
	log.Info("Updating data to database")
	for _, thread := range threads {
		sqlQuery := `UPDATE threads SET notified = $1 WHERE id = $2;`
		if _, err := a.DB.Exec(sqlQuery, true, thread.ID); err != nil {
			return err
		}
	}
	return nil
}
