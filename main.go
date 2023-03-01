package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

type App struct {
	ID          string
	Name        string
	Headline    string
	Description string
	URL         string
	ReviewCount int32
	Features    []string
	Categories  []string
}

type Review struct {
	AppID          string
	Rating         int32
	ReviewDate     time.Time
	Content        string
	ShopName       string
	Country        string
	TimeSpentOnApp time.Duration
}

func main() {
	seed := "https://apps.shopify.com"

	doc, err := getDoc(seed)
	if err != nil {
		panic(err)
	}

	fmt.Println(doc)
}

var urlChan = make(chan string, 1000)
var appChan = make(chan *App, 1000)
var reviewChan = make(chan *Review, 1000)

var httpClient *http.Client

func getHttpClient() *http.Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return httpClient
}

func request(url string) (*http.Response, error) {
	client := getHttpClient()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getDoc(url string) (*goquery.Document, error) {
	resp, err := request(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func ParseReviewList() {

}

func ParseReview() {

}
