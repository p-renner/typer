package main

import (
	"flag"
	"os"
)

var (
	quotesFile = flag.String("quotes", "quotes.json", "Path to quotes JSON file")
	QuoteID    = flag.Int("id", -1, "ID of the quote to use (default: random)")
	showHelp   = flag.Bool("help", false, "Show help message")
)

func Init() {
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}
}
