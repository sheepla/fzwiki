package client

import (
	"fmt"
	"net/url"
	"path"
)

// nolint:goerr113
func NewPageURL(title, lang string) string {
	u := &url.URL{
		Scheme: "https",
		Host: func(lang string) string {
			if lang == "" {
				return "en.wikipedia.org"
			}
			return fmt.Sprintf("%s.wikipedia.org", lang)
		}(lang),
		Path: path.Join("wiki", title),
	}
	return u.String()
}
