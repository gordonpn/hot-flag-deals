package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
)

type recipient struct {
	Name  string
	Email string
}

type thread struct {
	ID         int
	Title      string
	Link       string
	Posts      int
	Votes      int
	Views      int
	DatePosted time.Time
	Seen       bool
}

const (
	HCURL = "https://hc-ping.com"
)

func main() {
	_ = retrieveContent()
}

func init() {
	err := godotenv.Load()
	warnErr(err)
	log.SetLevel(log.DebugLevel)
}

func job() {
	signalHealthCheck("start")
	/*
		todo
	    send email
	    set those threads as seen
	*/
	signalHealthCheck("")
}

func signalHealthCheck(action string) {
	start, err := http.Get(fmt.Sprintf("%s/%s/%s", HCURL, os.Getenv("MAILER_HC_UUID"), action))
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Problem with GET request")
	}
	if start != nil {
		defer start.Body.Close()
	}
}

func warnErr(err error) {
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn()
	}
}

func retrieveContent() (threads []thread) {
	_, present := os.LookupEnv("DEV")
	host := "hotdeals_postgres"
	if present {
		host = "localhost"
	}
	port := 5432
	user := os.Getenv("POSTGRES_NONROOT_USER")
	password := os.Getenv("POSTGRES_NONROOT_PASSWORD")
	dbname := os.Getenv("POSTGRES_NONROOT_DB")

	pgURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", pgURI)
	if err != nil {
		log.Error("Error with opening connection with DB")
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Info("Successfully connected to DB")

	sqlStatement := `
  SELECT *
  FROM threads
  WHERE date_posted > CURRENT_TIMESTAMP - INTERVAL '30 day';`

	threadRows, err := db.Query(sqlStatement)
	warnErr(err)

	for threadRows.Next() {
		tempThread := thread{}
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

func filter(threads []thread) (filteredThreads []thread) {
	/*
	  todo
	    calculate both median and mean
	    if the two values are not too far apart
	    determine skewness
	    https://www.statisticshowto.com/pearsons-coefficient-of-skewness/
	    use the mean, otherwise use the median
	    filter out threads under that threshold and seen
	*/
	return
}

func sendMails(threads []thread, recipients []recipient) {
	//https://github.com/sendgrid/sendgrid-go/blob/master/examples/helpers/mail/example.go
	from := mail.NewEmail("Example User", "contact@gordon-pn.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", "gordon.pn6@gmail.com")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn()
	} else {
		log.WithFields(log.Fields{"Status Code": response.StatusCode}).Debug()
		log.WithFields(log.Fields{"Body": response.Body}).Debug()
		log.WithFields(log.Fields{"Headers": response.Headers}).Debug()
	}
}

func setSeen(threads []thread) {

}
