package main

import (
  "fmt"
  "github.com/gocolly/colly"
  "strings"
)

type thread struct {
  Id         string
  Title      string
  Link       string
  Posts      string
  Votes      string
  Views      string
  DatePosted string
}

// todo restructure this correctly
// todo connect to postgresql db to insert/update
// todo automate the scraper
// todo set up docker container

func main() {
}

func getPosts() (threads []thread) {
  collector := colly.NewCollector(
    colly.AllowedDomains("forums.redflagdeals.com"),
  )

  for i := 1; i <= 31; i++ {
    selector := fmt.Sprintf("#partition_forums > div > div.primary_content > div.forumbg > div > ul.topiclist.topics.with_categories > li:nth-child(%d)", i)
    collector.OnHTML(selector, func(element *colly.HTMLElement) {
      temp := thread{}
      temp.Id = strings.TrimSpace(element.Attr("data-thread-id"))
      titleSelector := "div > div.thread_info > div.thread_info_main.postvoting_enabled > div > h3 > a.topic_title_link"
      temp.Title = strings.TrimSpace(element.ChildText(titleSelector))
      temp.Link = fmt.Sprintf("%s/%s", "https://forums.redflagdeals.com/", strings.TrimSpace(element.ChildAttr(titleSelector, "href")))
      temp.Posts = strings.TrimSpace(element.ChildText("div > div.posts"))
      temp.Votes = strings.TrimSpace(element.ChildText("div > div.thread_info > div.thread_info_main.postvoting_enabled > div > div > dl > dd"))
      temp.Views = strings.TrimSpace(element.ChildText("div > div.views"))
      temp.DatePosted = strings.TrimSpace(element.ChildText("div > div.thread_info > div.thread_info_main.postvoting_enabled > div > div > div > span.first-post-time"))
      threads = append(threads, temp)
    })
  }

  collector.OnRequest(func(request *colly.Request) {
    fmt.Println("Visiting", request.URL.String())
  })

  for i := 1; i <= 10; i++ {
    url := fmt.Sprintf("https://forums.redflagdeals.com/hot-deals-f9/%d", i)
    collector.Visit(url)
  }

  return
}
