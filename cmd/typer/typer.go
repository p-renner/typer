package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"typer/quote"

	"golang.org/x/term"
)

var quotes quote.Quotes

func main() {
	Init()

	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		log.Fatal("Not running in a terminal")
	}

	if err := initQuotes(*quotesFile); err != nil {
		log.Fatal(err)
	}

	if err := typer(); err != nil {
		log.Fatal("Error:", err.Error())
	}
}

func initQuotes(path string) error {
	err := quotes.Load(path)

	if err != nil {
		return err
	}

	if quotes.Count() == 0 {
		return fmt.Errorf("no quotes found in %s", path)
	}

	return nil
}

const (
	CTRL_C = 3
	ESC    = 27
	BS     = 127
)

func typer() error {
	var err error
	var quote *quote.Quote

	if *QuoteID >= 0 {
		quote, err = quotes.GetByID(*QuoteID)
	} else {
		quote, err = quotes.GetRandom()
	}

	start := time.Now()

	defer func() {
		elapsed := time.Since(start)

		fmt.Printf("\nYou took: %s\n", elapsed.Truncate(time.Millisecond))

		if quote.Highscore == 0 {
			fmt.Printf("This was your first time, setting highscore to: %s\n", elapsed.Truncate(time.Millisecond))
			quote.Highscore = elapsed
			quotes.Save(*quotesFile)
			return
		}

		if elapsed > quote.Highscore {
			fmt.Printf("Your best time is: %s\n", quote.Highscore.Truncate(time.Millisecond))
			return
		}

		fmt.Printf("\nNew highscore! Previous best was: %s", quote.Highscore.Truncate(time.Millisecond))
		quote.Highscore = elapsed
		quotes.Save(*quotesFile)
	}()

	t := term.NewTerminal(os.Stdin, "")
	input := ""
	fmt.Fprint(t, quote.Quote)

	for {
		input, err = loop(input)

		// User requested exit
		if err != nil && err.Error() == "exiting" {
			return nil
		}

		if err != nil {
			return err
		}

		// TODO: if the quote is very long, the line will wrap and the clearing
		// will not work correctly. Might just clear the whole screen instead.
		clearLine()
		str := colorize(input, quote.Quote)
		fmt.Fprint(t, str)

		if input == quote.Quote {
			break
		}
	}

	if input == quote.Quote {
		fmt.Fprintf(t, "\nWell done!")
	}

	return nil
}

func colorize(input, quote string) string {
	result := ""
	correctLen := correctUpTo(input, quote)

	if correctLen > 0 {
		result += fmt.Sprintf("\033[32m%s\033[0m", input[:correctLen]) // Green
	}

	if correctLen < len(quote) {
		result += fmt.Sprintf("\033[31m%s\033[0m", quote[correctLen:len(input)]) // Red
	}

	if len(input) < len(quote) {
		result += quote[len(input):]
	}

	return result
}

func correctUpTo(input, quote string) int {
	correctLen := 0

	for i := 0; i < len(input) && i < len(quote); i++ {
		if input[i] != quote[i] {
			break
		}

		correctLen++
	}

	return correctLen
}

func loop(input string) (string, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		return "", err
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	key, err := readKeypress()

	if err != nil {
		return "", err
	}

	if key == CTRL_C || key == ESC {
		return input, fmt.Errorf("exiting")
	}

	if key == BS && len(input) > 0 {
		return input[:len(input)-1], nil
	}

	return input + string(key), nil
}

func clearLine() {
	fmt.Print("\r\033[K")
}

func readKeypress() (byte, error) {
	b := make([]byte, 1)
	n, err := os.Stdin.Read(b)

	if err != nil {
		return 0, err
	}

	if n == 0 {
		return 0, fmt.Errorf("no input read")
	}

	return b[0], nil
}
