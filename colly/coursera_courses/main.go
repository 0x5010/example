package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Course struct {
	Title       string
	Description string
	Creator     string
	Level       string
	URL         string
	Language    string
	Commitment  string
	HowToPass   string
	Rating      string
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("coursera.org", "www.coursera.org"),
		colly.CacheDir("./coursera_cache"),
	)

	detailCollector := c.Clone()

	courses := make([]Course, 0, 200)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Attr("class") == "Button_1qxkboh-o_O-primary_cv02ee-o_O-md_28awn8-o_O-primaryLink_109aggg" {
			return
		}
		link := e.Attr("href")
		if !strings.HasPrefix(link, "/browse") || strings.Index(link, "=signup") > -1 || strings.Index(link, "=login") > -1 {
			return
		}
		e.Request.Visit(link)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML(`a[name]`, func(e *colly.HTMLElement) {
		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Index(courseURL, "coursera.org/learn") != -1 {
			detailCollector.Visit(courseURL)
		}
	})

	detailCollector.OnHTML(`div[id=rendered-content]`, func(e *colly.HTMLElement) {
		log.Println("Course found", e.Request.URL)
		title := e.ChildText(".course-title")
		if title == "" {
			log.Println("No title found", e.Request.URL)
		}
		course := Course{
			Title:       title,
			URL:         e.Request.URL.String(),
			Description: e.ChildText("div.content"),
			Creator:     e.ChildText("div.creator-names > span"),
		}

		e.ForEach("table.basic-info-table tr", func(_ int, el *colly.HTMLElement) {
			switch el.ChildText("td:first-child") {
			case "Language":
				course.Language = el.ChildText("td:nth-child(2)")
			case "Level":
				course.Level = el.ChildText("td:nth-child(2)")
			case "Commitment":
				course.Commitment = el.ChildText("td:nth-child(2)")
			case "How To Pass":
				course.HowToPass = el.ChildText("td:nth-child(2)")
			case "User Ratings":
				course.Rating = el.ChildText("td:nth-child(2) div:nth-of-type(2)")
			}
		})
		courses = append(courses, course)
	})

	c.Visit("https://coursera.org/browse")

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	enc.Encode(courses)
}
