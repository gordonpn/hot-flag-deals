package main

import (
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go"
)

func (s *subscriber) sendConfirmEmail() error {
	if err := s.Validate(); err != nil {
		log.Error(fmt.Sprintf("Error with validating subscriber before sending confirmation email: %v", err))
		return errors.New("an error has occurred")
	}
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = getEmailBody(s)
	request.Body = Body
	response, err := sendgrid.API(request)
	statusCode := 0
	if err != nil {
		log.Errorf("Error with SendGrid: %v", err)
	} else {
		statusCode = response.StatusCode
		log.WithFields(log.Fields{"Status Code": statusCode}).Debug()
		log.WithFields(log.Fields{"Body": response.Body}).Debug()
	}
	return nil
}

func getEmailBody(s *subscriber) []byte {
	log.Debug("Preparing confirmation email")
	m := mail.NewV3Mail()

	address := "deals@gordon-pn.com"
	name := "Deals by gordonpn"
	e := mail.NewEmail(name, address)
	m.SetFrom(e)

	m.SetTemplateID(os.Getenv("SENDGRID_TEMPLATE_CONFIRM"))

	p := mail.NewPersonalization()

	log.WithFields(log.Fields{
		"Name":  s.Name,
		"Email": s.Email,
	}).Debug()

	to := mail.NewEmail(s.Name, s.Email)
	p.AddTos(to)

	dateNow := time.Now()
	date := fmt.Sprintf("%s %d, %d", dateNow.Month(), dateNow.Day(), dateNow.Year())

	unsubscribeLink := fmt.Sprintf("https://deals.gordon-pn.com/unsubscribe?email=%s", s.Email)
	confirmLink := fmt.Sprintf("https://deals.gordon-pn.com/confirm?email=%s", s.Email)

	p.SetDynamicTemplateData("date", date)
	p.SetDynamicTemplateData("unsubscribeLink", unsubscribeLink)
	p.SetDynamicTemplateData("confirmLink", confirmLink)
	p.SetDynamicTemplateData("name", s.Name)

	m.AddPersonalizations(p)
	return mail.GetRequestBody(m)
}
