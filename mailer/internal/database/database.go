package database

import (
	"database/sql"
	"fmt"
	types "github.com/gordonpn/hot-flag-deals/internal/data"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type pgdb struct {
	Database *sql.DB
}

var postgresDB pgdb

func GetDB() pgdb {
	if postgresDB.Database != nil {
		err := postgresDB.Database.Ping()
		if err == nil {
			return postgresDB
		}
	}
	postgresDB.Database = connectDB()
	return postgresDB
}

func RetrieveThreads() (threads []types.Thread) {
	pgDatabase := GetDB()
	db := pgDatabase.Database

	sqlStatement := `
  SELECT *
  FROM threads
  WHERE date_posted > CURRENT_TIMESTAMP - INTERVAL '30 day';`

	threadRows, err := db.Query(sqlStatement)
	warnErr(err)

	for threadRows.Next() {
		tempThread := types.Thread{}
		err = threadRows.Scan(
			&tempThread.ID,
			&tempThread.Title,
			&tempThread.Link,
			&tempThread.Posts,
			&tempThread.Votes,
			&tempThread.Views,
			&tempThread.DatePosted,
			&tempThread.Seen,
		)
		warnErr(err)
		threads = append(threads, tempThread)
	}
	log.WithFields(log.Fields{
		"len(threads)": len(threads),
		"cap(threads)": cap(threads)},
	).Debug("Length and capacity of threads")
	return
}

func CleanUp() {
	pgDatabase := GetDB()
	db := pgDatabase.Database
	log.Debug("Closing connection with DB")
	err := db.Close()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("Error with closing connection to DB")
	}
}

func SetSeen(threads []types.Thread) {
	pgDatabase := GetDB()
	db := pgDatabase.Database
	for _, thread := range threads {
		sqlStatement := `
    UPDATE threads
    SET seen = $1
    WHERE id = $2;`

		_, err := db.Exec(sqlStatement, true, thread.ID)
		warnErr(err)
	}
}

func connectDB() *sql.DB {
	host := "postgres"
	port := 5432
	user := os.Getenv("POSTGRES_NONROOT_USER")
	password := os.Getenv("POSTGRES_NONROOT_PASSWORD")
	dbname := os.Getenv("POSTGRES_NONROOT_DB")
	pgURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var db *sql.DB
	var err error
	log.Info("Attempting to connect to DB")
	for i := 1; i < 6; i++ {
		db, err = sql.Open("postgres", pgURI)
		if err != nil {
			log.Error("Error with opening connection with DB")
			panic(err)
		}

		err = db.Ping()
		if err == nil {
			break
		}
		retryWait := i * i
		log.Info(fmt.Sprintf("Connection attempt %d unsuccessful, retrying in %d seconds...", i, retryWait))
		time.Sleep(time.Duration(retryWait) * time.Second)
	}
	if db == nil {
		log.Fatal("Could not connect to DB")
	}

	log.Info("Successfully connected to DB")
	return db
}

func warnErr(err error) {
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn()
	}
}
