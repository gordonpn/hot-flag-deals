package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sort"
	"time"
)

func (a *App) notify(threads []thread) error {
	log.Info("Sending notifications: start")
	log.Info(fmt.Sprintf("Notifcations to send: %d", len(threads)))
	slackWebhook, present := os.LookupEnv("SLACK_NOTIFIER_HOOK")
	if !present {
		log.Panicf("Missing Slack Webhook URL")
	}
	sort.SliceStable(threads, func(i, j int) bool {
		return threads[i].Votes > threads[j].Votes
	})
	for _, thread := range threads {
		time.Sleep(5 * time.Second)
		reqBody, err := json.Marshal(map[string]string{
			"text": fmt.Sprintf("<%s|%s>", thread.Link, thread.Title),
		})
		if err != nil {
			return err
		}
		_, err = http.Post(slackWebhook, "application/json", bytes.NewBuffer(reqBody)) //nolint:gosec
		if err != nil {
			return err
		}
	}
	log.Info("Sending notifications: end")
	return nil
}
