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

type exitCode int

const (
	appVersion  = "0.0.7"
	appName     = "fzwiki"
	envNameLang = "FZWIKI_LANG"

	exitCodeOK exitCode = iota
	exitCodeErr
	exitCodeErrFuzzyFinder
	exitCodeErrWebBrowser
)

type options struct {
	Version  bool   `short:"V" long:"version" description:"Show version"`
	Open     bool   `short:"o" long:"open" description:"Open URL in your web browser"`
	Language string `short:"l" long:"lang" description:"Language for wikipedia.org such as \"en\", \"ja\", ..."`
}

func main() {
	os.Exit(int(Main(os.Args[1:])))
}

func Main(args []string) exitCode {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = appName
	parser.Usage = "[OPTIONS] QUERY..."
	args, err := parser.ParseArgs(args)
	if err != nil {
		if flags.WroteHelp(err) {
			return exitCodeOK
		} else {
			fmt.Fprintln(os.Stderr, "Argument parsing failed.")
			return exitCodeErr
		}
	}

	if opts.Version {
		fmt.Printf("%s: v%s\n", appName, appVersion)
		return exitCodeOK
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Must require argument(s).")
		return exitCodeErr
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
		return exitCodeErr
	}

	choices, err := find(result)
	if err != nil {
		log.Fatal(err)
		return exitCodeErrFuzzyFinder
	}

	for _, idx := range choices {
		url := createPageURL(result.Query.Search[idx].Title, lang)
		if opts.Open {
			if err := webbrowser.Open(url); err != nil {
				log.Fatal(err)
				return exitCodeErrWebBrowser
			}
		} else {
			fmt.Println(url)
		}
	}

	return exitCodeOK
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

func render(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		render(c, buf)
	}
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

func createPreview(i, w, h int, result *client.SearchResult) string {
	title := result.Query.Search[i].Title
	snippet := result.Query.Search[i].Snippet
	timestamp := result.Query.Search[i].Timestamp
	if s, err := html2text(snippet); err == nil {
		snippet = s
	}
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		title,
		runewidth.Wrap(snippet, w/2-5),
		humanize.Time(timestamp),
	)
}

func find(result *client.SearchResult) (choices []int, err error) {
	return fuzzyfinder.FindMulti(
		result.Query.Search,
		func(i int) string { return result.Query.Search[i].Title },
		fuzzyfinder.WithPreviewWindow(
			func(i, w, h int) string {
				if i == -1 {
					return ""
				}
				return createPreview(i, w, h, result)
			},
		),
	)
}
