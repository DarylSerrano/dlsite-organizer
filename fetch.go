package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// A Work contains information about a DLSite work
type Work struct {
	ID          string
	Name        string
	Circle      string
	VoiceActors []string
	Tags        []string
	sfw         bool
}

func parseName(goquerySelection *goquery.Selection) string {
	return goquerySelection.Find("#work_name > a").Text()
}

func parseCircle(goquerySelection *goquery.Selection) string {
	return goquerySelection.Find("#work_maker > tbody > tr > td > span > a").Text()
}

func parseVoiceActors(goquerySelection *goquery.Selection) []string {
	var voiceActors []string
	goquerySelection.
		Find("#work_outline > tbody").
		Find(`th:contains("声優")`).
		Next().
		Find("a").Each(func(_ int, s *goquery.Selection) {
		voiceActors = append(voiceActors, s.Text())
	})
	return voiceActors
}

func parseTags(goquerySelection *goquery.Selection) []string {
	var tags []string
	goquerySelection.
		Find("#work_outline > tbody").
		Find(`th:contains("ジャンル")`).
		Next().
		Find("a").Each(func(_ int, s *goquery.Selection) {
		tags = append(tags, s.Text())
	})
	return tags
}

func parseSfw(goquerySelection *goquery.Selection) bool {
	ageRestriction := goquerySelection.
		Find("#work_outline > tbody").
		Find(`th:contains("年齢指定")`).
		Next().
		Find("span").Text()

	return !strings.Contains(ageRestriction, "18")
}

func getUrlfromCode(code string) string {
	return fmt.Sprintf("https://www.dlsite.com/maniax/work/=/product_id/RJ%s.html", code)
}

func fetchWork(code string) (*Work, error) {
	url := getUrlfromCode(code)
	c := colly.NewCollector()
	var workInfo Work
	var err error

	workInfo.ID = code

	c.OnHTML("html", func(e *colly.HTMLElement) {
		goquerySelection := e.DOM

		workInfo.Name = parseName(goquerySelection)
		// fmt.Println(workInfo.Name)

		workInfo.Circle = parseCircle(goquerySelection)
		// fmt.Println(workInfo.Circle)

		workInfo.VoiceActors = parseVoiceActors(goquerySelection)
		// fmt.Println(workInfo.VoiceActors)

		workInfo.Tags = parseTags(goquerySelection)
		// fmt.Println(workInfo.Tags)

		workInfo.sfw = parseSfw(goquerySelection)

		//Find("#work_outline > tbody").

		// if ret, err := goquerySelection.
		// 	Find("#work_outline > tbody").
		// 	Find(`th:contains("ジャンル")`).
		// 	Next().
		// 	Find("a").
		// 	Html(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println(ret)
		// }
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, e error) {
		err = fmt.Errorf("Request URL: %s failed with response: %s \nError: %s", r.Request.URL.String(), r, e)
	})

	c.Visit(url)
	c.Wait()

	if err != nil {
		return nil, err
	}

	return &workInfo, nil
}
