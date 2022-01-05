package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	humanize "github.com/dustin/go-humanize"
	flags "github.com/jessevdk/go-flags"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-runewidth"
	"github.com/sheepla/fzwiki/client"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/net/html"
)

const (
	appVersion  = "0.0.5"
	appName     = "fzwiki"
	envNameLang = "FZWIKI_LANG"
)

type options struct {
	Version  bool   `short:"V" long:"version" description:"Show version"`
	Open     bool   `short:"o" long:"open" description:"Open URL in your web browser"`
	Language string `short:"l" long:"lang" description:"Language for wikipedia.org such as \"en\", \"ja\", ..."`
}

var opts options

func render(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		render(c, buf)
	}
}

func html2text(content string) (string, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	render(doc, &buf)
	return buf.String(), nil
}

func main() {
	os.Exit(Main(os.Args))
}

func Main(args []string) int {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = appName
	parser.Usage = "[OPTIONS] QUERY..."
	args, err := parser.ParseArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Argument parsing failed.")
		return 1
	}

	if opts.Version {
		fmt.Printf("%s: v%s\n", appName, appVersion)
		return 0
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Must require argument(s).")
		return 1
	}

	var lang string
	if opts.Language == "" {
		lang = os.Getenv(envNameLang)
	} else {
		lang = opts.Language
	}

	result, err := searchArticles(strings.Join(args, " "), lang)
	if err != nil {
		log.Fatal(err)
		return 1
	}
	for i := 0; i < len(result.Query.Search); i++ {
		if t, err := html2text(result.Query.Search[i].Title); err == nil {
			result.Query.Search[i].Title = t
		}
		if t, err := html2text(result.Query.Search[i].Snippet); err == nil {
			result.Query.Search[i].Snippet = t
		}
	}

	choices, err := fuzzyfinder.FindMulti(
		result.Query.Search,
		func(i int) string { return result.Query.Search[i].Title },
		fuzzyfinder.WithPreviewWindow(
			func(i, w, h int) string {
				if i == -1 {
					return ""
				}
				return fmt.Sprintf(
					"%s\n\n%s\n\n%s",
					result.Query.Search[i].Title,
					runewidth.Wrap(result.Query.Search[i].Snippet, w/2-5),
					humanize.Time(result.Query.Search[i].Timestamp),
				)
			},
		),
	)

	if err != nil {
		log.Fatal(err)
		return 2
	}

	for _, idx := range choices {
		url := createPageURL(result.Query.Search[idx].Title, lang)
		if opts.Open {
			if err := webbrowser.Open(url); err != nil {
				log.Fatal(err)
				return 3
			}
		} else {
			fmt.Println(url)
		}
	}

	return 0
}

func searchArticles(query, lang string) (*client.SearchResult, error) {
	url := client.CreateSearchURL(query, lang)
	return client.Execute(url)
}

func createPageURL(title, lang string) string {
	u := &url.URL{}
	u.Scheme = "https"
	if lang == "" {
		u.Host = fmt.Sprintf("%s.wikipedia.org", "en")
	} else {
		u.Host = fmt.Sprintf("%s.wikipedia.org", lang)
	}
	u.Path = fmt.Sprintf("wiki/%s", title)
	return u.String()
}
