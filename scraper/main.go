package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/whiteshtef/clockwork"
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

const (
	healthCheckURL = "https://hc-ping.com"
)

func main() {
	_, present := os.LookupEnv("DEV")
	if present {
		job()
	} else {
		scheduler := clockwork.NewScheduler()
		scheduler.SetPollingInterval(15 * 60000)
		scheduler.Schedule().Every(20).Minutes().Do(job)
		scheduler.Run()
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Problem with loading .env file")
	}
	log.SetLevel(log.DebugLevel)
}

func job() {
	signalHealthCheck("/start")

	threads := getPosts()
	upsertIntoDB(threads)

	signalHealthCheck("")
}

func signalHealthCheck(action string) {
	start, err := http.Get(fmt.Sprintf("%s/%s%s", healthCheckURL, os.Getenv("SCRAPER_HC_UUID"), action))
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Problem with GET request")
	}
	if start != nil {
		defer start.Body.Close()
	}
}

func getPosts() (threads []thread) {
	collector := colly.NewCollector(
		colly.AllowedDomains("forums.redflagdeals.com"),
	)

	titleSelector := "div > div.thread_info > div.thread_info_main.postvoting_enabled > div > h3"
	dateSelector := "div > div.thread_info > div.thread_info_main.postvoting_enabled > div > div > div > span.first-post-time"
	linkSelector := "div > div.thread_info > div.thread_info_main.postvoting_enabled > div > h3 > a.topic_title_link"
	linkPrefix := "https://forums.redflagdeals.com"

	for i := 1; i <= 31; i++ {
		selector := fmt.Sprintf("#partition_forums > div > div.primary_content > div.forumbg > div > ul.topiclist.topics.with_categories > li:nth-child(%d)", i)
		collector.OnHTML(selector, func(element *colly.HTMLElement) {
			tempThread := thread{}

			id := strToInt(element.Attr("data-thread-id"))
			if id == 0 {
				return
			}
			retailer := element.ChildText("div > div.thread_info > div.thread_info_main.postvoting_enabled > div > h3 > a.topictitle_retailer")
			posts := strToInt(element.ChildText("div > div.posts"))
			votes := strToInt(element.ChildText("div > div.thread_info > div.thread_info_main.postvoting_enabled > div > div > dl > dd"))
			views := strToInt(element.ChildText("div > div.views"))
			title := strings.TrimSpace(element.ChildText(titleSelector))
			title = strings.ReplaceAll(title, "\n", "")
			datePosted := strings.TrimSpace(element.ChildText(dateSelector))

			datetime := parseDateTime(datePosted)

			tempThread.ID = id
			if len(retailer) > 0 {
				tempThread.Title = fmt.Sprintf("[%s] %s", retailer, title)
			} else {
				tempThread.Title = title
			}
			tempThread.Link = fmt.Sprintf("%s%s", linkPrefix, strings.TrimSpace(element.ChildAttr(linkSelector, "href")))
			tempThread.Posts = posts
			tempThread.Votes = votes
			tempThread.Views = views
			tempThread.DatePosted = datetime
			tempThread.Seen = false
			tempThread.Notified = false

			threads = append(threads, tempThread)
		})
	}

	collector.OnRequest(func(request *colly.Request) {
		log.WithFields(log.Fields{"URL": request.URL.String()}).Info("Visiting")
	})

	for i := 1; i <= 10; i++ {
		url := fmt.Sprintf("https://forums.redflagdeals.com/hot-deals-f9/%d", i)
		err := collector.Visit(url)
		if err != nil {
			log.Warn(err)
		}
	}

	return
}

func upsertIntoDB(threads []thread) {
	host := "postgres"
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
	defer db.Close()

	log.WithFields(log.Fields{
		"len(threads)": len(threads),
		"cap(threads)": cap(threads)},
	).Debug("Length and capacity of threads")

	for _, thread := range threads {
		sqlStatement := `INSERT INTO threads (id, title, link, posts, votes, views, date_posted)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    ON CONFLICT (id)
    DO UPDATE SET title = EXCLUDED.title, posts = EXCLUDED.posts, votes = EXCLUDED.votes, views = EXCLUDED.views`

		_, err = db.Exec(sqlStatement, thread.ID, thread.Title, thread.Link, thread.Posts, thread.Votes, thread.Views, thread.DatePosted)
		if err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Problem with inserting")
		}
	}
}

func parseDateTime(datetime string) (parsedDateTime time.Time) {
	loc, _ := time.LoadLocation("America/Montreal")
	layout := "Jan 2 2006 3:04 pm"

	slices := strings.Fields(datetime)
	month := slices[0]
	dayOrdinal := slices[1]
	day := dayOrdinal[:len(dayOrdinal)-3]
	year := slices[2]
	hoursMinutes := slices[3]
	period := slices[4]

	datetimefmt := fmt.Sprintf("%s %s %s %s %s", month, day, year, hoursMinutes, period)
	parsedDateTime, err := time.ParseInLocation(layout, datetimefmt, loc)

	if err != nil {
		log.WithFields(log.Fields{
			"datetime":       datetime,
			"datetimefmt":    datetimefmt,
			"month":          month,
			"day":            day,
			"year":           year,
			"hoursMinutes":   hoursMinutes,
			"period":         period,
			"parsedDateTime": parsedDateTime,
		}).Debug("Parsing date and time")
		panic(err)
	}
	return
}

func strToInt(str string) (i int) {
	str = strings.TrimSpace(str)
	if len(str) < 1 {
		return 0
	}
	if strings.Contains(str, ",") {
		str = strings.Replace(str, ",", "", -1)
	}
	if strings.Contains(str, "+") {
		str = strings.Replace(str, "+", "", -1)
	}
	nonFractionalPart := strings.Split(str, ".")
	i, err := strconv.Atoi(nonFractionalPart[0])
	if err != nil {
		panic(err)
	}
	return
}
