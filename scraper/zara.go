package scraper

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type zara struct {
	log    logrus.FieldLogger
	client *http.Client
}

func NewZaraScraper(log logrus.FieldLogger) *zara {
	return &zara{
		log:    log,
		client: http.DefaultClient,
	}
}

func (z *zara) Scrape(task IScraperTask) {
	z.log.Debugf("scraping %s", task.GetSite())
}
