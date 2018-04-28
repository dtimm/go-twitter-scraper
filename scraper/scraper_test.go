package scraper_test

import (
	"github.com/dtimm/go-twitter-scraper/scraper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scraper", func() {
	Context("when called on a user", func() {
		It("returns all of their tweets", func() {
			scrape := scraper.New()
			tweets, err := scrape.User("vietarmis")
			Expect(err).ToNot(HaveOccurred())

			Expect(tweets).To(HaveLen(49))
		})
	})
})
