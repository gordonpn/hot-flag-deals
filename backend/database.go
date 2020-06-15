package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func connectDB(a *App, pgURI string) {
	var err error
	log.Info("Attempting to connect to DB")
	for i := 1; i < 6; i++ {
		a.DB, err = sql.Open("postgres", pgURI)
		if err != nil {
			log.Error("Error with opening connection with DB")
			panic(err)
		}

		err = a.DB.Ping()
		if err == nil {
			break
		}
		retryWait := i * i
		log.Info(fmt.Sprintf("Connection attempt %d unsuccessful, retrying in %d seconds...", i, retryWait))
		time.Sleep(time.Duration(retryWait) * time.Second)
	}
	if a.DB == nil {
		log.Fatal("Could not connect to DB")
	}
	log.Info("Successfully connected to DB")
}

func connectRedis(a *App) {
	host, exists := os.LookupEnv("REDIS_HOST")
	if !exists {
		host = "redis"
	}

	addr := fmt.Sprintf("%s:6379", host)
	var ctx = context.Background()
	a.RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	_, err := a.RDB.Ping(ctx).Result()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error with opening connection with Redis")
		panic(err)
	}
	log.Info("Successfully connected to Redis")
}
