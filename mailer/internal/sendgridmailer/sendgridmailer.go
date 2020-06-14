package sendgridmailer

import (
	"fmt"
	"github.com/gordonpn/hot-flag-deals/internal/data"
	"github.com/gordonpn/hot-flag-deals/internal/database"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func SendNewsletter(threads []types.Thread) {
	subscribers := getSubscribers()

	for _, subscriber := range subscribers {
		request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
		request.Method = "POST"
		var Body = getEmailBody(threads, subscriber)
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
	}
}

func getSubscribers() (subscribers []types.Subscriber) {
	pgDatabase := database.GetDB()
	db := pgDatabase.Database

	sqlStatement := `SELECT * FROM subscribers WHERE confirmed;`

	subscribersRow, err := db.Query(sqlStatement)
	warnErr(err)

	for subscribersRow.Next() {
		tempSub := types.Subscriber{}
		err = subscribersRow.Scan(
			&tempSub.ID,
			&tempSub.Name,
			&tempSub.Email,
			&tempSub.Confirmed,
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

func getEmailBody(threads []types.Thread, subscriber types.Subscriber) []byte {
	m := mail.NewV3Mail()

	address := "deals@gordon-pn.com"
	name := "Deals by gordonpn"
	e := mail.NewEmail(name, address)
	m.SetFrom(e)

	m.SetTemplateID(os.Getenv("SENDGRID_TEMPLATE"))

	p := mail.NewPersonalization()

	log.WithFields(log.Fields{
		"Name":  subscriber.Name,
		"Email": subscriber.Email,
	}).Debug()

	to := mail.NewEmail(subscriber.Name, subscriber.Email)
	p.AddTos(to)

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

func warnErr(err error) {
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn()
	}
}
