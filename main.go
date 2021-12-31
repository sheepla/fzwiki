package main

import (
	"fmt"
	"log"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/sheepla/fzwiki/client"
	"github.com/toqueteos/webbrowser"
	"os"
)

type Options struct {
	Open     bool   `short:"o" long:"open" description:"open URL with webbrowser"`
	Language string `short:"l" long:"lang" description:"wikipedia language"`
}

var opts Options

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "fzwiki"
	parser.Usage = "[OPTIONS] QUERY..."
	args, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if len(args) == 0 {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	result := searchArticles(strings.Join(args, " "), opts.Language)

	choices, err := fuzzyfinder.FindMulti(
		result.Query.Search,
		func(i int) string { return result.Query.Search[i].Title },
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("%s\n\n%s", result.Query.Search[i].Title, result.Query.Search[i].Snippet)
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
    if lang == "" {
        return fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", "en", title)
    } else {
        return fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", lang, title)
    }
}
