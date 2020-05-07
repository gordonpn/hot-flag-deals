package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
)

type recipient struct {
}
type thread struct {
}

const (
	HCURL = "https://hc-ping.com"
)

func main() {

}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn()
	}
	log.SetLevel(log.DebugLevel)
}

func job() {
	signalHealthCheck("start")
	/*
			  get all threads from DB from the last month
			  calculate both median and mean
			  if the two values are not too far apart
		    determine skewness
		    https://www.statisticshowto.com/pearsons-coefficient-of-skewness/
			  use the mean, otherwise use the median
			  filter out threads under that threshold and seen
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

func retrieveContent() []thread {

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
