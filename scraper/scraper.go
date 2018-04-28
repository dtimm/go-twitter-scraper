package scraper

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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

	ioutil.WriteFile("tmp.json", []byte(body), 0644)
	return parseBody(body)
}

func parseBody(body string) ([]string, error) {
	tweets := twitterBody{}

	json.Unmarshal([]byte(body), &tweets)

	fmt.Print(tweets)

	return []string{}, nil
}

type twitterBody struct {
	Items []tweet `json:"items_html"`
	Min   string  `json:"min_position"`
	More  bool    `json:"has_more_items"`
}

type tweet struct {
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	TweetID    string `json:"tweet_ids"`
	Text       string `json:"emojified_text_as_html"`
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

	htmlClean := html.UnescapeString(string(b))
	utf8Clean := strings.Replace(htmlClean, `\u007b`, "{", -1)
	utf8Clean = strings.Replace(utf8Clean, `\u007d`, "}", -1)
	// utf8Clean = strings.Replace(utf8Clean, `\n`, "\n", -1)
	utf8Clean = strings.Replace(utf8Clean, `\u003c`, "<", -1)
	utf8Clean = strings.Replace(utf8Clean, `\u003e`, ">", -1)

	return utf8Clean, nil
}
