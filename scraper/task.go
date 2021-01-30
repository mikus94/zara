package scraper

import "net/url"

type IScraperTask interface {
	GetSite() *url.URL
	GetSizes() []string
}
