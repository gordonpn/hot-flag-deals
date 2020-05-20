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
	"github.com/whiteshtef/clockwork"
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
	_, present := os.LookupEnv("DEV")
	if present {
		job()
	} else {
		scheduler := clockwork.NewScheduler()
		scheduler.SetPollingInterval(30 * 60000)
		scheduler.Schedule().Every().Day().At("10:00").Do(job)
		scheduler.Run()
	}
}

func init() {
	err := godotenv.Load()
	warnErr(err)
	log.SetLevel(log.DebugLevel)
}

func job() {
	signalHealthCheck("/start")

	threads := retrieveThreads()
	filteredThreads := filter(threads)
	if len(filteredThreads) > 0 {
		ok := sendNewsletter(filteredThreads)
		if ok {
			setSeen(filteredThreads)
		}
	}
	cleanUp()

	signalHealthCheck("")
}

func signalHealthCheck(action string) {
	start, err := http.Get(fmt.Sprintf("%s/%s%s", HCURL, os.Getenv("MAILER_HC_UUID"), action))
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
	host := "hotdeals_postgres"
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
	postgresDB := app{}
	postgresDB.Database = db
	return postgresDB
}

func retrieveThreads() (threads []thread) {
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
		standDevThreshold   = 0.9
		viewsMean           float64
		viewsMedian         int
		viewsSkewness       float64
		viewsSlice          []int
		viewsStandDev       float64
		viewsThresholdCoeff = 1.05
		votesMean           float64
		votesMedian         int
		votesSkewness       float64
		votesSlice          []int
		votesStandDev       float64
		votesThresholdCoeff = 5.0
	)
	for _, thread := range threads {
		viewsSlice = append(viewsSlice, thread.Views)
		votesSlice = append(votesSlice, thread.Votes)
	}
	viewsMean = GetMean(viewsSlice)
	votesMean = GetMean(votesSlice)

	viewsMedian = GetMedian(viewsSlice)
	votesMedian = GetMedian(votesSlice)

	viewsStandDev = GetStandDev(viewsSlice, viewsMean)
	votesStandDev = GetStandDev(votesSlice, votesMean)

	viewsSkewness = GetSkewness(viewsMean, viewsMedian, viewsStandDev)
	votesSkewness = GetSkewness(votesMean, votesMedian, votesStandDev)

	if math.Abs(viewsSkewness) >= standDevThreshold {
		viewsThreshold = viewsMedian
	} else {
		viewsThreshold = Round(viewsMean)
	}
	if math.Abs(votesSkewness) >= standDevThreshold {
		votesThreshold = votesMedian
	} else {
		votesThreshold = Round(votesMean)
	}
	viewsThreshold = Round(float64(viewsThreshold) * viewsThresholdCoeff)
	votesThreshold = Round(float64(votesThreshold) * votesThresholdCoeff)
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

func GetMean(intSlice []int) (mean float64) {
	sum := 0
	for _, num := range intSlice {
		sum += num
	}
	mean = float64(sum) / float64(len(intSlice))
	return
}

func GetMedian(intSlice []int) (median int) {
	sort.Ints(intSlice)
	middle := len(intSlice) / 2
	if len(intSlice)%2 == 0 {
		median = (intSlice[middle-1] + intSlice[middle]) / 2
	} else {
		median = intSlice[middle]
	}
	return
}

func GetStandDev(intSlice []int, mean float64) (standDev float64) {
	for i := range intSlice {
		standDev += math.Pow(float64(intSlice[i])-mean, 2)
	}
	standDev = math.Sqrt(standDev / float64(len(intSlice)))
	return
}

func GetSkewness(mean float64, median int, standDev float64) (skewness float64) {
	skewness = (mean - float64(median)) * 3 / standDev
	return
}

func filter(threads []thread) (filteredThreads []thread) {
	viewsThreshold, votesThreshold := getThresholds(threads)
	const TimeThreshold = 72

	for _, thread := range threads {
		if (thread.Views >= viewsThreshold || thread.Votes >= votesThreshold) && !thread.Seen {
			timeNow := time.Now()
			diffHours := timeNow.Sub(thread.DatePosted).Hours()
			if diffHours <= TimeThreshold {
				filteredThreads = append(filteredThreads, thread)
			}
		}
	}
	sort.SliceStable(filteredThreads, func(this, that int) bool {
		return filteredThreads[this].Votes > filteredThreads[that].Votes
	})
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
			&tempSub.Name,
			&tempSub.Email,
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
		log.WithFields(log.Fields{
			"Name":  subscriber.Name,
			"Email": subscriber.Email,
		}).Debug()
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

func sendNewsletter(threads []thread) bool {
	// todo fix: one 'to' per email
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = getEmailBody(threads)
	request.Body = Body
	response, err := sendgrid.API(request)
	statusCode := 0
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn()
	} else {
		statusCode = response.StatusCode
		log.WithFields(log.Fields{"Status Code": statusCode}).Debug()
		log.WithFields(log.Fields{"Body": response.Body}).Debug()
	}
	return statusCode >= 200 && statusCode < 300
}

func setSeen(threads []thread) {
	db := connectDB().Database
	for _, thread := range threads {
		sqlStatement := `
    UPDATE threads
    SET seen = $1
    WHERE id = $2;`

		_, err := db.Exec(sqlStatement, true, thread.ID)
		warnErr(err)
	}
}

func cleanUp() {
	db := connectDB().Database
	log.Debug("Closing connection with DB")
	err := db.Close()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("Error with closing connection to DB")
	}
}

func Round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}
