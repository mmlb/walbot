package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

var catRe = regexp.MustCompile("/category/([a-z]+)")

func parseFrontPage() map[string]bool {
	fmt.Fprintln(os.Stderr, "dongerlist front page parsing started")
	defer fmt.Fprintln(os.Stderr, "dongerlist front page parsing done")

	// Create a collector
	c := colly.NewCollector()

	categories := map[string]bool{}
	c.OnHTML("option", func(e *colly.HTMLElement) {
		// fmt.Println(strings.ToLower(e.Text), e.Attr("value"))
		value := e.Attr("value")
		if !catRe.MatchString(value) {
			return
		}
		if found := categories[value]; found {
			return
		}
		categories[value] = true
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping
	c.Visit("http://dongerlist.com")
	return categories
}

func setupDongerlist(w *walbot) {
	c := colly.NewCollector()

	allExcluderRe := regexp.MustCompile("^/page/")
	c.OnHTML("a.nextpostslink", func(e *colly.HTMLElement) {
		if allExcluderRe.MatchString(e.Request.URL.Path) {
			return
		}
		c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	dongers := map[string][]string{}
	c.OnHTML("textarea.donger", func(e *colly.HTMLElement) {
		parts := catRe.FindStringSubmatch(e.Request.URL.Path)
		if parts == nil {
			return
		}
		cat := parts[1]
		dongers[cat] = append(dongers[cat], e.Text)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	fmt.Fprintln(os.Stderr, "dongerlist categories parsing started")
	for cat := range parseFrontPage() {
		c.Visit(cat)
	}
	fmt.Fprintln(os.Stderr, "dongerlist categories parsing done")

	for cat, dongers := range dongers {
		thinge := "donger-" + cat
		_, err := w.makeThinge(thinge)
		if err != nil && !errors.Is(err, errThingeExists) {
			fmt.Fprintf(os.Stderr, "error saving donger: %v\n", err)
			continue
		}
		w.Store.Set("thinge."+thinge, dongers)
	}
}
