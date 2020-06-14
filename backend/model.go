package main

import (
	"database/sql"
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
}

type subscriber struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Confirmed bool   `json:"confirmed"`
}

func getThreads(db *sql.DB) ([]thread, error) {
	log.Debug("Deals threads requested from database")
	sqlStatement := `SELECT * FROM threads WHERE date_posted > CURRENT_TIMESTAMP - INTERVAL '2 day' AND votes > 0;`

	rows, err := db.Query(sqlStatement)

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
		); err != nil {
			return nil, err
		}
		threads = append(threads, tempThread)
	}
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
		return err
	}
	return nil
}

func (s *subscriber) deleteSubscriber(db *sql.DB) error {
	log.Debug("Deleting subscriber")
	// TODO validate
	sqlQuery := `DELETE FROM subscribers WHERE email = $1`
	_, err := db.Exec(sqlQuery, s.Email)

	if err != nil {
		return err
	}
	return nil
}

func (s *subscriber) updateSubscriber(db *sql.DB) error {
	log.Debug("Updating subscriber")
	// TODO validate
	sqlQuery := `UPDATE subscribers SET confirmed = TRUE WHERE email = $1`
	_, err := db.Exec(sqlQuery, s.Email)

	if err != nil {
		return err
	}
	return nil
}
