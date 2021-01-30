package main

import (
	"os"
	"time"

	"github.com/mikus94/zara/reader"
	"github.com/mikus94/zara/scraper"
	"github.com/mikus94/zara/structs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DEFAULT_EXECUTION_TIME = 6 * time.Hour
	DEFAULT_CHECK_TIMEOUT  = 5 * time.Minute
)

func main() {
	logger := logrus.New()
	// logger.SetLevel(logrus.DebugLevel)

	path := "./url_sites.csv"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		logger.WithError(err).Fatal("cannot find configuration file")
	}

	execution_time := viper.GetDuration("execution_time")
	if execution_time.Seconds() < 0 || execution_time.Hours() > 6 {
		logger.Warnf(
			"execution time passed was: %s, using default %s",
			execution_time, DEFAULT_EXECUTION_TIME,
		)
		execution_time = DEFAULT_EXECUTION_TIME
	}
	check_timeout := viper.GetDuration("check_timeout")
	if check_timeout.Seconds() < 1 || check_timeout.Hours() > 1 {
		logger.Warnf(
			"check timeout passed was: %s, using default %s",
			check_timeout, DEFAULT_CHECK_TIMEOUT,
		)
		check_timeout = DEFAULT_CHECK_TIMEOUT
	}

	logger.Infof("Using path: %s", path)
	// END OF READING CONFIGS

	//////////////////////// INITIALIZE SCRAP OBJECTS //////////////////////////
	lines, err := reader.ReadCSVFile(path)
	if err != nil {
		logger.WithError(err).Fatal("cannot read csv file")
	}
	objectsToScrap := make([]*structs.ScrapObject, 0, len(lines))
	for i, l := range lines {
		if i == 0 {
			// skip header
			continue
		}
		obj, err := structs.NewScrapObjectByCSVLine(l)
		if err != nil {
			logger.WithError(err).Error("cannot create scrap object")
			continue
		}
		objectsToScrap = append(objectsToScrap, obj)
	}
	////////////////////////////////////////////////////////////////////////////

	////////////////////////// INITIALIZE SCRAPERS /////////////////////////////
	zaraScraper := scraper.NewZaraScraper(logger)
	////////////////////////////////////////////////////////////////////////////

	////////////////////////////// ACTUAL RUN //////////////////////////////////
	// initial check
	for _, task := range objectsToScrap {
		zaraScraper.Scrape(task)
	}
	for {
		select {
		case <-time.After(execution_time):
			// finish job
			logger.Info("finishing my job")
			break
		case <-time.After(check_timeout):
			// checking
			for _, task := range objectsToScrap {
				zaraScraper.Scrape(task)
			}
		}
	}
}
