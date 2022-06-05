package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
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

// nolint:gochecknoglobals
var (
	appVersion  = "unknown"
	appRevision = "unknown"
	appName     = "fzwiki"
	envNameLang = "FZWIKI_LANG"
)

const (
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
	code, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(int(code))
}

func run(args []string) (exitCode, error) {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = appName
	parser.Usage = "[OPTIONS] QUERY..."
	args, err := parser.ParseArgs(args)
	if err != nil {
		if flags.WroteHelp(err) {
			return exitCodeOK, nil
		} else {
			return exitCodeErr, fmt.Errorf("argument parsing failed: %w", err)
		}
	}

	if opts.Version {
		fmt.Printf("%s: v%s-rev%s\n", appName, appVersion, appRevision)

		return exitCodeOK, nil
	}

	if len(args) == 0 {
		return exitCodeErr, errors.New("must require arguments.")
	}

	var lang string
	if opts.Language == "" {
		lang = os.Getenv(envNameLang)
	} else {
		lang = opts.Language
	}

	result, err := searchArticles(strings.Join(args, " "), lang)
	if err != nil {
		return exitCodeErr, fmt.Errorf("failed to search articles: %w", err)
	}

	choices, err := find(result)
	if err != nil {
		if errors.Is(fuzzyfinder.ErrAbort, err) {
			return exitCodeOK, nil
		}
		return exitCodeErr, fmt.Errorf("an error occurred on fuzzyfinder: %w", err)
	}

	for _, idx := range choices {
		url := createPageURL(result.Query.Search[idx].Title, lang)
		if opts.Open {
			if err := webbrowser.Open(url); err != nil {
				return exitCodeErrWebBrowser, fmt.Errorf("an error occurred on opening web browser: %w", err)
			}
		} else {
			fmt.Println(url)
		}
	}

	return exitCodeOK, nil
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
	if lang == "" {
		lang = "en"
	}
	u := &url.URL{
		Scheme: "https",
		Path:   path.Join("wiki", title),
		Host:   fmt.Sprintf("%s.wikipedia.org", lang),
	}
	return u.String()
}

func createPreview(i, w, h int, result *client.SearchResult) string {
	title := result.Query.Search[i].Title
	snippet := result.Query.Search[i].Snippet
	timestamp := result.Query.Search[i].Timestamp
	wordcount := result.Query.Search[i].Wordcount
	if s, err := html2text(snippet); err == nil {
		snippet = s
	}
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s, %s words",
		title,
		runewidth.Wrap(snippet, w/2-5),
		humanize.Time(timestamp),
		humanize.Comma(wordcount),
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
