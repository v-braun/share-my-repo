package strategy

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/v-braun/hero-scrape"
)

var _ heroscrape.Strategy = (*GHStrategy)(nil)

var og = heroscrape.NewOgStrategy()
var hr = heroscrape.NewHeuristicStrategy()

type GHStrategy struct {
	OgResult  *heroscrape.SearchResult
	HrResult  *heroscrape.SearchResult
	EndResult *heroscrape.SearchResult
}

// NewOgStrategy returns a new Strategy that search for OG meta tags
func NewGitHubStrategy() *GHStrategy {
	return new(GHStrategy)
}

func isErrNotIncomplete(err error) bool {
	return err != nil && err != heroscrape.ErrNotComplete
}

func (gh *GHStrategy) Scrape(srcURL *url.URL, doc *goquery.Document) (*heroscrape.SearchResult, error) {
	ogRes, err := og.Scrape(srcURL, doc)
	if isErrNotIncomplete(err) {
		return nil, err
	}

	hrRes, err := hr.Scrape(srcURL, doc)
	if isErrNotIncomplete(err) {
		return nil, err
	}

	gh.OgResult = ogRes
	gh.HrResult = hrRes

	result := new(heroscrape.SearchResult)
	if ogRes != nil {
		result.Image = ogRes.Image
		result.Title = ogRes.Title
		result.Description = ogRes.Description
	}
	if hrRes != nil && hrRes.Image != "" {
		result.Image = hrRes.Image
	}

	gh.EndResult = result
	return result, nil
}
