package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/jessevdk/go-flags"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-runewidth"
	"github.com/sheepla/fzwiki/client"
	"github.com/toqueteos/webbrowser"
)

// nolint:gochecknoglobals
var (
	appName     = "fzwiki"
	appVerion   = "0.0.0"
	appRevision = "?????"
	appUsage    = "[OPTIONS] QUERY..."
)

const envNameLang = "FZWIKI_LANG"

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrSearch
	exitCodeErrFuzzyFinder
	exitCodeErrWebBrowser
)

type options struct {
	Version bool `short:"V" long:"version" description:"Show version"`
	Open    bool `short:"o" long:"open" description:"Open pages URL on the web browser"`
	// Preview bool `short:"p" long:"preview" description:"Preview page on the terminal"`
	LimitNum int    `short:"n" long:"num" description:"Max number of search items" default:"20"`
	Lang     string `short:"l" long:"language" description:"Language of wikipedia" default:"en"`
}

func main() {
	code, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(int(code))
}

// nolint:funlen,cyclop
func run(cliArgs []string) (exitCode, error) {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = appName
	parser.Usage = appUsage

	args, err := parser.ParseArgs(cliArgs)
	if err != nil {
		if flags.WroteHelp(err) {
			return exitCodeOK, nil
		}

		return exitCodeErrArgs, fmt.Errorf("argument parsing failed %w", err)
	}

	if opts.Version {
		// nolint:forbidigo
		fmt.Printf("%s: v%s-rev%s\n", appName, appVerion, appRevision)

		return exitCodeOK, nil
	}

	if len(args) == 0 {
		// nolint:goerr113
		return exitCodeErrArgs, errors.New("must reuire arguments")
	}

	lang := func() string {
		if envvar := strings.TrimSpace(os.Getenv(envNameLang)); envvar != "" {
			return envvar
		}

		return opts.Lang
	}()

	result, err := client.Search(client.Param{
		Query: strings.Join(args, " "),
		Limit: opts.LimitNum,
		Lang:  lang,
	})
	if err != nil {
		return exitCodeErrSearch, fmt.Errorf("failed to search articles:%w", err)
	}

	if result == nil {
		// nolint:goerr113
		return exitCodeErrSearch, errors.New("failed to search articles (result is nil)")
	}

	choises, err := findMulti(result)
	if err != nil {
		return exitCodeErrFuzzyFinder, fmt.Errorf("an error occurred on fuzzyfinder: %w", err)
	}

	for _, idx := range choises {
		url := client.NewPageURL(
			result.Query.Search[idx].Title,
			lang,
		)

		if opts.Open {
			if err := webbrowser.Open(url); err != nil {
				return exitCodeErrWebBrowser, fmt.Errorf("failed to open the URL (%s): %w", url, err)
			}
		} else {
			// nolint:forbidigo
			fmt.Println(url)

			return exitCodeOK, nil
		}
	}

	return exitCodeOK, nil
}

// func find(result *client.Result) (int, error) {
// 	return fuzzyfinder.Find(
// 		result.Query.Search,
// 		func(idx int) string {
// 			if idx <= 0 {
// 				return ""
// 			}
//
// 			return result.Query.Search[idx].Title
// 		},
// 		fuzzyfinder.WithPreviewWindow(func(idx, width, height int) string {
// 			if idx <= 0 {
// 				return ""
// 			}
// 			padding := 5
//
// 			return runewidth.Wrap(renderPreviewWindow(result, idx), width/2-padding)
// 		}),
// 	)
// }

// nolint:wrapcheck
func findMulti(result *client.Result) ([]int, error) {
	return fuzzyfinder.FindMulti(
		result.Query.Search,
		func(idx int) string {
			if idx == -1 {
				return ""
			}

			return result.Query.Search[idx].Title
		},
		fuzzyfinder.WithPreviewWindow(func(idx, width, height int) string {
			if idx == -1 {
				return ""
			}
			padding := 5

			return runewidth.Wrap(renderPreviewWindow(result, idx), width/2-padding)
		}),
	)
}

func renderPreviewWindow(result *client.Result, idx int) string {
	title := result.Query.Search[idx].Title
	timestamp := humanize.Time(result.Query.Search[idx].Timestamp)
	snippet := result.Query.Search[idx].Snippet
	words := result.Query.Search[idx].Wordcount

	return fmt.Sprintf(
		"%s [%s]\n\n%s\n\n%s",
		title,
		timestamp,
		snippet,
		sprintfUnlessEmpty("\n%d words", words),
	)
}

func sprintfUnlessEmpty(format string, a any) string {
	if format == "" || a == "" {
		return ""
	}

	return fmt.Sprintf(format, a)
}
