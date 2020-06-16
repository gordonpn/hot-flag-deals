package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/whiteshtef/clockwork"
	"os"
	"time"
)

type thread struct {
	ID         int
	Title      string
	Link       string
	Posts      int
	Votes      int
	Views      int
	DatePosted time.Time
	Seen       bool
	Notified   bool
}

type App struct {
	DB          *sql.DB
	threads     []thread
	votesMedian int
	postsMedian int
	viewsMedian int
}

func main() {
	_, present := os.LookupEnv("DEV")
	if present {
		job()
	} else {
		scheduler := clockwork.NewScheduler()
		scheduler.SetPollingInterval(20 * 60000)
		scheduler.Schedule().Every(30).Minutes().Do(job)
		scheduler.Run()
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Warnf("Did not load .env file: %v", err)
	}
	log.SetLevel(log.DebugLevel)
}

func job() {
	log.Info("Job start")
	signalHealthCheck("/start")

	a := App{}
	if err := a.connectDB(); err != nil {
		signalHealthCheck("/fail")
		log.Panicf("Error with connecting to database: %v", err)
	}
	defer a.DB.Close()
	if err := a.getThreads(); err != nil {
		signalHealthCheck("/fail")
		log.Panicf("Error with fetching data from database: %v", err)
	}
	a.getThresholds()
	threads := a.filter()
	if err := a.notify(threads); err != nil {
		signalHealthCheck("/fail")
		log.Panicf("Error with notifying: %v", err)
	}
	if err := a.markNotified(threads); err != nil {
		signalHealthCheck("/fail")
		log.Panicf("Error with updating data: %v", err)
	}

	signalHealthCheck("")
	log.Info("Job end")
}
