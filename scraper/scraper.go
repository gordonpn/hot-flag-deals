package main

import (
  "fmt"
  "github.com/gocolly/colly"
  "strconv"
  "strings"
)

type thread struct {
  Id         int
  Title      string
  Link       string
  Posts      int
  Votes      int
  Views      int
  DatePosted string
}

// todo parse date to a type compatible with postgresql
// todo connect to postgresql db to insert/update
// todo automate the scraper
// todo set up docker container

func main() {
  getPosts()
}

func getPosts() (threads []thread) {
  collector := colly.NewCollector(
    colly.AllowedDomains("forums.redflagdeals.com"),
  )

  for i := 1; i <= 31; i++ {
    selector := fmt.Sprintf("#partition_forums > div > div.primary_content > div.forumbg > div > ul.topiclist.topics.with_categories > li:nth-child(%d)", i)
    collector.OnHTML(selector, func(element *colly.HTMLElement) {
      temp := thread{}

      titleSelector := "div > div.thread_info > div.thread_info_main.postvoting_enabled > div > h3"
      dateSelector := "div > div.thread_info > div.thread_info_main.postvoting_enabled > div > div > div > span.first-post-time"
      linkSelector := "div > div.thread_info > div.thread_info_main.postvoting_enabled > div > h3 > a.topic_title_link"
      linkPrefix := "https://forums.redflagdeals.com"

      retailer := element.ChildText("div > div.thread_info > div.thread_info_main.postvoting_enabled > div > h3 > a.topictitle_retailer")
      id := StrToInt(element.Attr("data-thread-id"))
      posts := StrToInt(element.ChildText("div > div.posts"))
      votes := StrToInt(element.ChildText("div > div.thread_info > div.thread_info_main.postvoting_enabled > div > div > dl > dd"))
      views := StrToInt(element.ChildText("div > div.views"))
      title := strings.TrimSpace(element.ChildText(titleSelector))
      title = strings.ReplaceAll(title, "\n", "")

      temp.Id = id
      if len(retailer) > 0 {
        temp.Title = fmt.Sprintf("[%s] %s", retailer, title)
      } else {
        temp.Title = title
      }
      temp.Link = fmt.Sprintf("%s%s", linkPrefix, strings.TrimSpace(element.ChildAttr(linkSelector, "href")))
      temp.Posts = posts
      temp.Votes = votes
      temp.Views = views
      temp.DatePosted = strings.TrimSpace(element.ChildText(dateSelector))
      threads = append(threads, temp)
    })
  }

  collector.OnRequest(func(request *colly.Request) {
    fmt.Println("Visiting", request.URL.String())
  })

  for i := 1; i <= 1; i++ {
    url := fmt.Sprintf("https://forums.redflagdeals.com/hot-deals-f9/%d", i)
    collector.Visit(url)
  }

  fmt.Println(threads)

  return
}

func StrToInt(str string) (i int) {
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
