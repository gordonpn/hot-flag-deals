package main

import (
	"database/sql"
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

func getThreads(db *sql.DB) ([]thread, error) {
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
