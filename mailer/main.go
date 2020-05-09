package main

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
)

type subscriber struct {
	ID    int
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

type app struct {
	Database *sql.DB
}

const (
	HCURL = "https://hc-ping.com"
)

func main() {
	threads := retrieveContent()
	filteredThreads := filter(threads)
	sendNewsletter(filteredThreads)
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

func connectDB() app {
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
	postgresDB := app{}
	postgresDB.Database = db
	return postgresDB
}

func retrieveContent() (threads []thread) {
	db := connectDB().Database

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

func getThresholds(threads []thread) (viewsThreshold, votesThreshold int) {
	var (
		middle            int
		viewsMean         int
		viewsMedian       int
		viewsSkewness     float64
		viewsSlice        []int
		viewsStandDev     float64
		viewsSum          = 0
		votesMean         int
		votesMedian       int
		votesSkewness     float64
		votesSlice        []int
		votesStandDev     float64
		votesSum          = 0
		standDevThreshold = 0.9
		thresholdFactor   = 1.3
	)
	for _, thread := range threads {
		viewsSum += thread.Views
		viewsSlice = append(viewsSlice, thread.Views)
		votesSum += thread.Votes
		votesSlice = append(votesSlice, thread.Votes)
	}
	viewsMean = viewsSum / len(threads)
	votesMean = votesSum / len(threads)
	if len(threads)%2 == 0 {
		middle = len(threads) / 2
	} else {
		middle = (len(threads) - 1) / 2
	}
	sort.Ints(viewsSlice)
	sort.Ints(votesSlice)
	viewsMedian = viewsSlice[middle]
	votesMedian = votesSlice[middle]

	for i := range threads {
		viewsStandDev += math.Pow(float64(viewsSlice[i]-viewsMean), 2)
		votesStandDev += math.Pow(float64(votesSlice[i]-votesMean), 2)
	}
	viewsStandDev = math.Sqrt(viewsStandDev / float64(len(viewsSlice)))
	votesStandDev = math.Sqrt(votesStandDev / float64(len(votesSlice)))
	viewsSkewness = float64((viewsMean-viewsMedian)*3) / (viewsStandDev)
	votesSkewness = float64((votesMean-votesMedian)*3) / (votesStandDev)
	if math.Abs(viewsSkewness) >= standDevThreshold {
		viewsThreshold = viewsMedian
	} else {
		viewsThreshold = viewsMean
	}
	if math.Abs(votesSkewness) >= standDevThreshold {
		votesThreshold = votesMedian
	} else {
		votesThreshold = votesMean
	}
	viewsThreshold = round(float64(viewsThreshold) * thresholdFactor)
	votesThreshold = round(float64(votesThreshold) * thresholdFactor)
	log.WithFields(log.Fields{
		"viewsMean":              viewsMean,
		"viewsMedian":            viewsMedian,
		"viewsStandardDeviation": viewsStandDev,
		"viewsSkewness":          viewsSkewness,
		"viewsThreshold":         viewsThreshold,
	}).Debug()
	log.WithFields(log.Fields{
		"votesMean":              votesMean,
		"votesMedian":            votesMedian,
		"votesStandardDeviation": votesStandDev,
		"votesSkewness":          votesSkewness,
		"votesThreshold":         votesThreshold,
	}).Debug()
	return
}

func filter(threads []thread) (filteredThreads []thread) {
	viewsThreshold, votesThreshold := getThresholds(threads)

	for _, thread := range threads {
		if (thread.Views >= viewsThreshold && thread.Votes >= votesThreshold) && !thread.Seen {
			filteredThreads = append(filteredThreads, thread)
		}
	}
	log.WithFields(log.Fields{
		"len(filteredThreads)": len(filteredThreads),
		"cap(filteredThreads)": cap(filteredThreads)},
	).Debug("Length and capacity of filtered threads")
	return
}

func getSubscribers() (subscribers []subscriber) {
	db := connectDB().Database

	sqlStatement := `
  SELECT *
  FROM subscribers;`

	subscribersRow, err := db.Query(sqlStatement)
	warnErr(err)

	for subscribersRow.Next() {
		tempSub := subscriber{}
		err = subscribersRow.Scan(
			&tempSub.ID,
			&tempSub.Email,
			&tempSub.Name,
		)
		warnErr(err)
		subscribers = append(subscribers, tempSub)
	}
	log.WithFields(log.Fields{
		"len(subscribers)": len(subscribers),
		"cap(subscribers)": cap(subscribers)},
	).Debug("Length and capacity of subscribers")
	return
}

func getEmailBody(threads []thread) []byte {
	m := mail.NewV3Mail()

	address := "deals@gordon-pn.com"
	name := "Deals by gordonpn"
	e := mail.NewEmail(name, address)
	m.SetFrom(e)

	m.SetTemplateID(os.Getenv("SENDGRID_TEMPLATE"))

	p := mail.NewPersonalization()
	var tos []*mail.Email
	subscribers := getSubscribers()

	for _, subscriber := range subscribers {
		tos = append(tos, mail.NewEmail(subscriber.Name, subscriber.Email))
	}

	p.AddTos(tos...)

	dateNow := time.Now()
	date := fmt.Sprintf("%s %d, %d", dateNow.Month(), dateNow.Day(), dateNow.Year())

	p.SetDynamicTemplateData("date", date)

	var dealList []map[string]string
	var deal map[string]string

	for _, v := range threads {
		deal = make(map[string]string)
		deal["title"] = v.Title
		deal["link"] = v.Link
		dealList = append(dealList, deal)
	}

	p.SetDynamicTemplateData("deals", dealList)

	m.AddPersonalizations(p)
	return mail.GetRequestBody(m)
}

func sendNewsletter(threads []thread) {
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = getEmailBody(threads)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn()
	} else {
		log.WithFields(log.Fields{"Status Code": response.StatusCode}).Debug()
	}
}

func setSeen(threads []thread) {

}

func round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}
