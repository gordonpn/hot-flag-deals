package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gordonpn/hot-flag-deals/mailer/pkg/database"
	"github.com/gordonpn/hot-flag-deals/mailer/pkg/filter"

	"github.com/gordonpn/hot-flag-deals/mailer/internal/sendgridmailer"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/whiteshtef/clockwork"
)

const (
	healthCheckURL = "https://hc-ping.com"
)

func main() {
	_, present := os.LookupEnv("DEV")
	if !present {
		scheduler := clockwork.NewScheduler()
		scheduler.SetPollingInterval(30 * 60000)
		scheduler.Schedule().Every().Day().At("8:30").Do(job)
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

	threads := database.RetrieveThreads()
	filteredThreads := filter.Filter(threads)
	if len(filteredThreads) > 0 {
		sendgridmailer.SendNewsletter(filteredThreads)
		database.SetSeen(filteredThreads)
	}
	database.CleanUp()

	signalHealthCheck("")
}

func signalHealthCheck(action string) {
	resp, err := http.Get(fmt.Sprintf("%s/%s%s", healthCheckURL, os.Getenv("MAILER_HC_UUID"), action))
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Problem with GET request")
	}
	if resp != nil {
		defer resp.Body.Close()
	}
}

func warnErr(err error) {
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn()
	}
}
