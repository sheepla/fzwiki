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
    AppVersion = "0.0.1"
    AppName = "fzwiki"
)

type Options struct {
    Version  bool   `short:"V" long:"version" description:"show version"`
	Open     bool   `short:"o" long:"open" description:"open URL in your web browser"`
	Language string `short:"l" long:"lang" description:"language for wikipedia.org such as \"en\", \"ja\", ..."`
}

var opts Options

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
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = AppName
	parser.Usage = "[OPTIONS] QUERY..."
	args, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

    if opts.Version {
        fmt.Printf("%s: v%s\n", AppName, AppVersion)
        os.Exit(0)
    }

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Must require argument(s).")
		os.Exit(1)
	}

	result := searchArticles(strings.Join(args, " "), opts.Language)
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
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf(
				"%s\n\n%s\n\n%s",
				result.Query.Search[i].Title,
				runewidth.Wrap(result.Query.Search[i].Snippet, w),
				humanize.Time(result.Query.Search[i].Timestamp),
			)
		},
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	for _, i := range choices {
		url := createPageUrl(result.Query.Search[i].Title, opts.Language)
		if opts.Open {
			if err := webbrowser.Open(url); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println(url)
		}
	}
}

func searchArticles(query, lang string) client.SearchResult {
	url := client.CreateUrl(query, lang)
	result := client.Execute(url)
	return result
}

func createPageUrl(title, lang string) string {
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
