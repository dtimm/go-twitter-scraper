package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	htmlParse "golang.org/x/net/html"
)

type scraper struct {
	client  http.Client
	headers map[string]string
}

// New creates a default scraper instance
func New() scraper {
	return scraper{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		headers: map[string]string{
			"Accept":                "application/json, text/javascript, */*; q=0.01",
			"User-Agent":            "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36",
			"X-Twitter-Active-User": "yes",
			"X-Requested-With":      "XMLHttpRequest",
		},
	}
}

// User gets all the tweets of a Twitter user and returns them.
func (s scraper) User(name string) ([]string, error) {
	s.headers["Referer"] = fmt.Sprintf("https://twitter.com/%s", name)

	body, err := s.getBody(name)
	if err != nil {
		return nil, err
	}

	return parseBody(body)
}

func parseBody(body string) ([]string, error) {
	tweets := twitterBody{}

	json.Unmarshal([]byte(body), &tweets)

	ioutil.WriteFile("tmp.html", []byte(tweets.Items), 0644)

	node, err := htmlParse.Parse(bytes.NewReader([]byte(tweets.Items)))
	if err != nil {
		return nil, err
	}

	tweetText := make([]string, 0, 50)

	parseNodes(node, 0)

	return tweetText, nil
}

func parseNodes(n *htmlParse.Node, d int) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {

		parseNodes(c, d+1)
	}
}

type twitterBody struct {
	Items  string `json:"items_html"`
	MinPos string `json:"min_position"`
	MaxPos string `json:"max_position"`
	More   bool   `json:"has_more_items"`
}

func (s scraper) getBody(name string) (string, error) {
	s.headers["Referer"] = fmt.Sprintf("https://twitter.com/%s", name)
	address := fmt.Sprintf("https://twitter.com/i/profiles/show/%s/timeline/tweets?include_available_features=1&include_entities=1&include_new_items_bar=true", name)
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return "", err
	}

	for key, val := range s.headers {
		req.Header.Add(key, val)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
