package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type App struct {
	Name        string
	Headline    string
	Description string
	URL         string
	ReviewCount string
	Features    []string
	Categories  []string
}

var skipPages = []string{
	"/login",
	"/signup",
	"/reviews",
	"/partners",
}

func main() {

	var apps []*App
	var appIndex map[string]*App = make(map[string]*App)

	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}

		// write to json file

		content, err := json.Marshal(apps)
		if err != nil {
			fmt.Println(err)
			return
		}

		os.WriteFile("apps.json", content, 0644)
	}()

	c := colly.NewCollector(
		colly.Async(),
		colly.AllowedDomains("apps.shopify.com"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
		RandomDelay: 5 * time.Second,
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		target := e.Attr("href")

		if target == "" {
			return
		}

		if _, exist := appIndex[target]; exist {
			return
		}

		parsedUrl, _ := url.Parse(target)
		for _, skipPage := range skipPages {
			if strings.Contains(parsedUrl.Path, skipPage) {
				return
			}
		}

		if len(apps) < 100 {
			e.Request.Visit(target)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("body[id=AppDetailsShow]", func(e *colly.HTMLElement) {

		fmt.Println("Process app: ", e.Request.URL)

		app, err := parseApp(e.DOM)
		if err != nil {
			fmt.Println(err)
			return
		}

		// remove query string
		app.URL = strings.Split(e.Request.URL.String(), "?")[0]

		apps = append(apps, app)
		appIndex[e.Request.URL.String()] = app

	})

	c.Visit("https://apps.shopify.com")

	signChan := make(chan os.Signal, 1)
	go func() {
		c.Wait()
		signChan <- syscall.SIGTERM
	}()

	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	<-signChan

	fmt.Println("Done")
}

func parseApp(dom *goquery.Selection) (*App, error) {
	var app App

	app.Name = strings.TrimSpace(dom.Find("#adp-hero h1").First().Text())
	app.Headline = strings.TrimSpace(dom.Find("#app-details > h2").First().Text())
	app.Description = strings.TrimSpace(dom.Find("#app-details > p").First().Text())
	app.ReviewCount = strings.TrimSpace(dom.Find("a[href=\"#adp-reviews\"]").First().Text())

	dom.Find("#app-details > ul > li").Each(func(i int, s *goquery.Selection) {
		app.Features = append(app.Features, strings.TrimSpace(s.Text()))
	})

	dom.Find("#adp-details-section a[href^=\"https://apps.shopify.com/categories\"]").Each(func(i int, s *goquery.Selection) {
		app.Categories = append(app.Categories, strings.TrimSpace(s.Text()))
	})

	return &app, nil
}
