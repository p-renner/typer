package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"typer/quote"
)

var (
	quotesFile = flag.String("quotes", "quotes.json", "Path to quotes JSON file")
	QuoteID    = flag.Int("id", -1, "ID of the quote to use (default: random)")
	showHelp   = flag.Bool("help", false, "Show help message")
)

func Init() quote.Quotes {
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	quotes, err := initQuotes(*quotesFile)

	if err != nil {
		log.Fatal(err)
	}

	return quotes
}

func initQuotes(path string) (quote.Quotes, error) {
	var quotes quote.Quotes
	err := quotes.Load(path)

	if err != nil {
		return nil, err
	}

	if quotes.Count() == 0 {
		return nil, fmt.Errorf("no quotes found in %s", path)
	}

	return quotes, nil
}
