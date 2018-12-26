package strategy

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/v-braun/hero-scrape"
)

var _ heroscrape.Strategy = (*ghStrategy)(nil)

var og = heroscrape.NewOgStrategy()
var hr = heroscrape.NewHeuristicStrategy()

type ghStrategy struct {
}

// NewOgStrategy returns a new Strategy that search for OG meta tags
func NewGitHubStrategy() heroscrape.Strategy {
	return new(ghStrategy)
}

func isErrNotIncomplete(err error) bool {
	return err != nil && err != heroscrape.ErrNotComplete
}

func (gh *ghStrategy) Scrape(srcURL *url.URL, doc *goquery.Document) (*heroscrape.SearchResult, error) {
	ogRes, err := og.Scrape(srcURL, doc)
	if isErrNotIncomplete(err) {
		return nil, err
	}

	hrRes, err := hr.Scrape(srcURL, doc)
	if isErrNotIncomplete(err) {
		return nil, err
	}

	result := new(heroscrape.SearchResult)
	if ogRes != nil {
		result.Image = ogRes.Image
		result.Title = ogRes.Title
		result.Description = ogRes.Description
	}
	if hrRes != nil && hrRes.Image != "" {
		result.Image = hrRes.Image
	}

	return result, nil
}
