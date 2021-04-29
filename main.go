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

func main() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("error loading location 'Asia/Bangkok': %v\n", err)
	}

	time.Local = ict
	c := colly.NewCollector(colly.AllowedDomains("www.set.or.th"))

	// calendar
	c.OnHTML("a[data-target='#calendarDetailModal']", func(e *colly.HTMLElement) {
		time.Sleep(100 * time.Millisecond)
		onclick := e.Attr("onclick")
		r := regexp.MustCompile(`symbol\=(\w*?)\&`)
		strMatch := r.FindStringSubmatch(onclick)
		if len(strMatch) == 0 {
			return
		}
		stock := strMatch[1]

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
		fmt.Print(dt.Format("2006-01-02"), "     ", strings.TrimSpace(stock))

		link := fmt.Sprintf("https://www.set.or.th/set/companyprofile.do?symbol=%s&ssoPageId=4&language=en&country=US", stock)
		if err := c.Visit(e.Request.AbsoluteURL(link)); err != nil {
			fmt.Println("Error visit:", link, err.Error())
		}

	})

	// Dvd. Yield(%)
	c.OnHTML("#maincontent > div > div.table-reponsive > table > tbody > tr:nth-child(2) > td > div:nth-child(2) > div.col-xs-12.col-md-8 > div.col-xs-9.col-md-5", func(e *colly.HTMLElement) {
		fmt.Printf("     %s\n", strings.TrimSpace(e.Text))
	})

	if err := c.Visit("https://www.set.or.th/set/xcalendar.do?eventType=XD&index=2&language=en&country=US"); err != nil {
		log.Fatalln(err.Error())
	}
}
