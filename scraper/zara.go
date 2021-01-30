package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/mikus94/zara/notifier"
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

const (
	productDetailsBoxSelector = ".product-detail-info"
	ulSizeListSelector        = "ul.product-size-selector__size-list"

	sizeDisabledSelector   = "product-size-selector__size-list-item--is-disabled"
	sizeOutOfStockSelector = "product-size-selector__size-list-item--out-of-stock"
	sizeComingSoonSelector = "product-size-selector__size-list-item--back-soon"
)

func checkIfDisabled(s *goquery.Selection) bool {
	// element that we check should be li type
	if !s.Is("li") {
		return true
	}
	return s.HasClass(sizeDisabledSelector) ||
		s.HasClass(sizeOutOfStockSelector) ||
		s.HasClass(sizeComingSoonSelector)
}

func (z *zara) Scrape(task IScraperTask) {
	logger := z.log.WithFields(logrus.Fields{
		"method":   "Scrape",
		"scraper":  "Zara",
		"taskSite": task.GetSite().String(),
	})
	logger.Debugf("scraping %s", task.GetSite())
	res, err := z.client.Get(task.GetSite().String())
	if err != nil {
		logger.WithError(err).Error("cannot get site")
		return
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Fatal(err)
	}
	doc.Find(productDetailsBoxSelector).Each(func(i int, s *goquery.Selection) {
		productName := s.Find("h1").First().Text()
		logger = logger.WithField("product", productName)
		sizeList := s.Find(ulSizeListSelector).First()
		sizeList.Find("li").Each(func(i int, s *goquery.Selection) {
			// we get li element of size list
			size := s.Text()
			logger = logger.WithField("size", size)
			if checkIfDisabled(s) {
				logger.Debug("size not available")
				return
			}
			logger.Debug("available")
			notifier.Notify(logger, &notifier.NotificationMessage{
				Title:    "Zara has your size",
				Subtitle: productName,
				Message:  fmt.Sprintf("%s is available! Hurry up!", size),
				Url:      task.GetSite().String(),
				Sound:    gosxnotifier.Default,
				Hash:     fmt.Sprintf("%s.%s", strings.ReplaceAll(productName, " ", "-"), size),
			})
		})
	})
	logger.Debug("processing")
}
