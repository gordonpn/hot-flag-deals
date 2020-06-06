package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Problem with loading .env file")
	}
	log.SetLevel(log.DebugLevel)
}

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("POSTGRES_NONROOT_USER"),
		os.Getenv("POSTGRES_NONROOT_PASSWORD"),
		os.Getenv("POSTGRES_NONROOT_DB"))

	a.Run(":8080")
}
