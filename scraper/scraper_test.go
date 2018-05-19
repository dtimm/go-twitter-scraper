package scraper_test

import (
	"github.com/dtimm/go-twitter-scraper/scraper"

	"testing"

	. "github.com/onsi/gomega"
)

func TestScraper(t *testing.T) {
	g := NewGomegaWithT(t)
	scrape := scraper.New()
	tweets, err := scrape.User("vietarmis")
	g.Expect(err).ToNot(HaveOccurred())

	g.Expect(tweets).To(HaveLen(49))
}
