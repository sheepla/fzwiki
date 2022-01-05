package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type SearchResult struct {
	Query struct {
		Search []struct {
			Title     string    `json:"title"`
			Pageid    int       `json:"pageid"`
			Snippet   string    `json:"snippet"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"search"`
	} `json:"query"`
}

func CreateSearchURL(query, lang string) string {
	u := &url.URL{}
	u.Scheme = "https"
	if lang == "" {
		u.Host = fmt.Sprintf("%s.wikipedia.org", "en")
	} else {
		u.Host = fmt.Sprintf("%s.wikipedia.org", lang)
	}
	u.Path = "w/api.php"

	q := u.Query()
	q.Set("action", "query")
	q.Set("format", "json")
	q.Set("list", "search")
	q.Set("srsearch", query)
	u.RawQuery = q.Encode()
	return u.String()
}

func Execute(url string) (*SearchResult, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
