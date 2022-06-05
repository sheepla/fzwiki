package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Result struct {
	Query struct {
		Search []struct {
			Title     string    `json:"title"`
			Pageid    int64     `json:"pageid"`
			Snippet   string    `json:"snippet"`
			Timestamp time.Time `json:"timestamp"`
			Wordcount int64     `json:"wordcount"`
		} `json:"search"`
	} `json:"query"`
}

type Param struct {
	Query string
	Limit int
	Lang  string
}

// nolint:exhaustivestruct,exhaustruct,varnamelen
func newURL(param Param) string {
	u := &url.URL{
		Scheme: "https",
		Host: func(lang string) string {
			if lang == "" {
				return "en.wikipedia.org"
			}

			return fmt.Sprintf("%s.wikipedia.org", lang)
		}(param.Lang),
		Path: "w/api.php",
	}

	val := u.Query()
	val.Set("action", "query")
	val.Set("format", "json")
	val.Set("list", "search")
	val.Set("srsearch", param.Query)
	val.Set("srlimit", strconv.Itoa(param.Limit))
	u.RawQuery = val.Encode()

	return u.String()
}

// nolint:goerr113
func Search(param Param) (*Result, error) {
	res, err := http.Get(newURL(param))
	if err != nil {
		return nil, fmt.Errorf("failed to search articles: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search articles (status code: %d, status:%s)", res.StatusCode, res.Status)
	}

	var result Result
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	return &result, nil
}
