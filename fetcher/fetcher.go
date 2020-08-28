package fetcher

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

var reCode = regexp.MustCompile(`RG\d+`)

// A Work contains information about a DLSite work
type Work struct {
	ID          string
	Name        string
	Circle      CircleParsed
	VoiceActors []string
	Tags        []string
	SFW         bool
}

// A CircleParsed contains information of DLSite group parsed
type CircleParsed struct {
	ID   string
	Name string
}

func getRGCode(filename string) string {
	foundRj := reCode.FindString(filename)

	return foundRj[2:]
}

func parseName(goquerySelection *goquery.Selection) string {
	return goquerySelection.Find("#work_name > a").Text()
}

func parseCircle(goquerySelection *goquery.Selection) CircleParsed {
	var data CircleParsed
	aTag := goquerySelection.Find("#work_maker > tbody > tr > td > span > a")
	data.Name = aTag.Text()
	url, urlExists := aTag.Attr("href")
	if urlExists {
		data.ID = getRGCode(url)
	}

	return data
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

// FetchWork scraps information of the work from DLSite
func FetchWork(code string) (*Work, error) {
	url := getUrlfromCode(code)
	c := colly.NewCollector()
	var workInfo Work
	var err error

	workInfo.ID = code

	c.OnHTML("html", func(e *colly.HTMLElement) {
		goquerySelection := e.DOM

		workInfo.Name = parseName(goquerySelection)

		workInfo.Circle = parseCircle(goquerySelection)

		workInfo.VoiceActors = parseVoiceActors(goquerySelection)

		workInfo.Tags = parseTags(goquerySelection)

		workInfo.SFW = parseSfw(goquerySelection)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, e error) {
		err = fmt.Errorf("Request URL: %v failed with response: %v \nError: %s", r.Request.URL.String(), r, e)
	})

	c.Visit(url)
	c.Wait()

	if err != nil {
		return nil, err
	}

	return &workInfo, nil
}
