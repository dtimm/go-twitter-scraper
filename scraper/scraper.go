package scraper

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"time"
)

type scraper struct {
	client  http.Client
	headers map[string]string
}

func New() scraper {
	return scraper{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		headers: map[string]string{
			"Accept":                "application/json, text/javascript, */*; q=0.01",
			"Referer":               "https://twitter.com/%s",
			"User-Agent":            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
			"X-Twitter-Active-User": "yes",
			"X-Requested-With":      "XMLHttpRequest",
		},
	}
}

// User gets all the tweets of a Twitter user and returns them.
func (s scraper) User(name string) ([]string, error) {
	address := fmt.Sprintf("https://twitter.com/i/profiles/show/%s/timeline/tweets?include_available_features=1&include_entities=1&include_new_items_bar=true", name)
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range s.headers {
		req.Header.Add(key, val)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ioutil.WriteFile("tmp.json", []byte(html.UnescapeString(string(b))), 0644)
	return []string{}, nil
}
