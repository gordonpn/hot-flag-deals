package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"time"
)

type thread struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Link       string    `json:"link"`
	Posts      int       `json:"posts"`
	Votes      int       `json:"votes"`
	Views      int       `json:"views"`
	DatePosted time.Time `json:"date"`
	Seen       bool      `json:"seen"`
	Notified   bool      `json:"notified"`
}

type subscriber struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Confirmed bool   `json:"confirmed"`
}

func getThreads(a *App) ([]thread, error) {
	log.Debug("Deals threads requested")
	var ctx = context.Background()
	redisResults, err := a.RDB.Get(ctx, "threads").Result()
	if err != redis.Nil {
		var results []thread
		err := json.Unmarshal([]byte(redisResults), &results)
		if err == nil {
			log.Debug("Returning data from Redis")
			return results, nil
		}
	}
	sqlStatement := `SELECT * FROM threads WHERE date_posted > CURRENT_TIMESTAMP - INTERVAL '2 day' AND votes > 0;`

	rows, err := a.DB.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var threads []thread
	for rows.Next() {
		tempThread := thread{}
		if err = rows.Scan(
			&tempThread.ID,
			&tempThread.Title,
			&tempThread.Link,
			&tempThread.Posts,
			&tempThread.Votes,
			&tempThread.Views,
			&tempThread.DatePosted,
			&tempThread.Seen,
			&tempThread.Notified,
		); err != nil {
			return nil, err
		}
		threads = append(threads, tempThread)
	}
	redisThreads, err := json.Marshal(threads)
	if err != nil {
		log.Error(err)
	}
	_, err = a.RDB.SetNX(ctx, "threads", redisThreads, time.Hour).Result()
	if err != nil {
		log.Error(err)
	}
	log.Debug("Returning data from database")
	return threads, nil
}

func (s *subscriber) createSubscriber(db *sql.DB) error {
	log.Debug("Creating subscriber")
	// TODO validate name and email
	sqlQuery := `INSERT INTO subscribers (name, email)
    VALUES ($1, $2)
    ON CONFLICT (email)
    DO UPDATE SET confirmed = TRUE`

	_, err := db.Exec(sqlQuery, s.Name, s.Email)

	if err != nil {
		log.Error(err)
		return errors.New("an error has occurred")
	}
	return nil
}

func (s *subscriber) deleteSubscriber(db *sql.DB) error {
	log.Debug("Deleting subscriber")
	// TODO validate email
	sqlQuery := `DELETE FROM subscribers WHERE email = $1`
	_, err := db.Exec(sqlQuery, s.Email)

	if err != nil {
		log.Error(err)
		return errors.New("an error has occurred")
	}
	return nil
}

func (s *subscriber) updateSubscriber(db *sql.DB) error {
	log.Debug("Updating subscriber")
	// TODO validate email
	sqlQuery := `UPDATE subscribers SET confirmed = TRUE WHERE email = $1`
	_, err := db.Exec(sqlQuery, s.Email)

	if err != nil {
		log.Error(err)
		return errors.New("an error has occurred")
	}
	return nil
}
