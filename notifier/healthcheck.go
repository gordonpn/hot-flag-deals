package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func signalHealthCheck(action string) {
	HCURL := "https://hc-ping.com"
	if _, present := os.LookupEnv("DEV"); present {
		return
	}
	log.Info(fmt.Sprintf("Signalling Healthchecks.io: %s", action))
	uuid, present := os.LookupEnv("NOTIFIER_HC_UUID")
	if !present {
		log.Panicf("Running in production and missing health check UUID for notifier")
	}
	req, err := http.Get(fmt.Sprintf("%s/%s%s", HCURL, uuid, action))
	if err != nil {
		log.Warnf("Error with healthcheck signal: %v", err)
	}
	if req != nil {
		defer req.Body.Close()
	}
}
