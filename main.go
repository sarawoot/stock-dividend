package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Stock struct {
	Name   string
	XDDate string
	Yield  float64
	Set    bool
}

func main() {
	stocks := map[string]*Stock{}
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("error loading location 'Asia/Bangkok': %v\n", err)
	}

	time.Local = ict
	c := colly.NewCollector(colly.AllowedDomains("www.set.or.th", "marketdata.set.or.th"))

	// calendar
	c.OnHTML("a[data-target='#calendarDetailModal']", func(e *colly.HTMLElement) {
		time.Sleep(100 * time.Millisecond)
		onclick := e.Attr("onclick")
		r := regexp.MustCompile(`symbol\=(\w*?)\&`)
		strMatch := r.FindStringSubmatch(onclick)
		if len(strMatch) == 0 {
			return
		}
		stockName := strings.TrimSpace(strMatch[1])

		r = regexp.MustCompile(`xDate\=(\d*?)\&`)
		strMatch = r.FindStringSubmatch(onclick)
		if len(strMatch) == 0 {
			return
		}
		millisecStr := strMatch[1]
		millisec, err := strconv.Atoi(millisecStr)
		if err != nil {
			return
		}

		dt := time.Unix(0, int64(millisec)*1000000)

		stocks[stockName] = &Stock{
			Name:   stockName,
			XDDate: dt.Format("2006-01-02"),
		}

		link := fmt.Sprintf("https://www.set.or.th/set/companyprofile.do?symbol=%s&ssoPageId=4&language=en&country=US", stockName)
		if err := c.Visit(e.Request.AbsoluteURL(link)); err != nil {
			fmt.Println("Error visit:", link, err.Error())
		}

	})

	// Dvd. Yield(%)
	c.OnHTML("#maincontent > div > div.table-reponsive > table > tbody > tr:nth-child(2) > td > div:nth-child(2) > div.col-xs-12.col-md-8 > div.col-xs-9.col-md-5", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		r := regexp.MustCompile(`symbol\=(\w*?)\&`)
		strMatch := r.FindStringSubmatch(url)
		if len(strMatch) == 0 {
			return
		}
		stockName := strings.TrimSpace(strMatch[1])
		yield, err := strconv.ParseFloat(strings.TrimSpace(e.Text), 64)
		if err != nil {
			return
		}
		s := stocks[stockName]
		if s != nil {
			s.Yield = yield
		}
	})

	c.OnHTML("#maincontent > div > div:nth-child(4) > div > div > div > div:nth-child(7) > table > tbody", func(e *colly.HTMLElement) {
		hrefs := e.ChildAttrs("tr > td:nth-child(1) > a", "href")
		for _, v := range hrefs {
			r := regexp.MustCompile(`symbol\=(\w*?)\&`)
			strMatch := r.FindStringSubmatch(v)
			if len(strMatch) == 0 {
				return
			}
			stockName := strings.TrimSpace(strMatch[1])
			s := stocks[stockName]
			if s != nil {
				s.Set = true
			}
		}
	})

	if err := c.Visit("https://www.set.or.th/set/xcalendar.do?eventType=XD&index=2&language=en&country=US"); err != nil {
		log.Fatalln(err.Error())
	}

	// if err := c.Visit("https://marketdata.set.or.th/mkt/sectorquotation.do?sector=SET100&language=en&country=US"); err != nil {
	if err := c.Visit("https://marketdata.set.or.th/mkt/sectorquotation.do?sector=SETHD&language=en&country=US"); err != nil {
		fmt.Println("set100 error: ", err.Error())
	}

	for _, v := range stocks {
		if v.Set {
			fmt.Println(*v)
		}
	}
}
