package scraper

type IScraper interface {
	// Scrape is executing scraping task and returning list
	// of available sizes or error while scraping
	// if no size was found then return empty array.
	Scrape(task IScraperTask) ([]string, error)
}
