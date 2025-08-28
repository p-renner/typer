package quote_test

import (
	"os"
	"testing"
	"typer/quote"
)

func TestStruct(t *testing.T) {
	q := quote.Quote{
		Quote:  "Test Quote",
		Author: "Tester",
	}

	if q.Quote != "Test Quote" {
		t.Errorf("Expected 'Test Quote', got '%s'", q.Quote)
	}

	if q.Author != "Tester" {
		t.Errorf("Expected 'Tester', got '%s'", q.Author)
	}
}

func TestLoad(t *testing.T) {
	jsonData := `[{"quote":"Test Quote","author":"Tester"}]`
	fileName := "test_quotes.json"

	err := os.WriteFile(fileName, []byte(jsonData), 0644)

	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	defer os.Remove(fileName)

	var quotes quote.Quotes
	err = quotes.Load(fileName)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if quotes.Count() != 1 {
		t.Errorf("Expected count 1, got %d", quotes.Count())
	}

	if quotes[0].Quote != "Test Quote" || quotes[0].Author != "Tester" {
		t.Errorf("Quote data mismatch")
	}
}

func TestCount(t *testing.T) {
	quotes := quote.Quotes{
		{Quote: "Quote 1", Author: "Author 1"},
		{Quote: "Quote 2", Author: "Author 2"},
	}

	if quotes.Count() != 2 {
		t.Errorf("Expected count 2, got %d", quotes.Count())
	}
}

func TestGetByID(t *testing.T) {
	quotes := quote.Quotes{
		{Quote: "Quote 1", Author: "Author 1"},
		{Quote: "Quote 2", Author: "Author 2"},
	}

	quote, err := quotes.GetByID(1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if quote == nil || quote.Quote != "Quote 2" {
		t.Errorf("Expected 'Quote 2', got '%v'", quote)
	}

	quote, err = quotes.GetByID(5)

	if err != nil {
		t.Errorf("Expected no error for out-of-bounds, got %v", err)
	}

	if quote != nil {
		t.Errorf("Expected nil for out-of-bounds, got '%v'", quote)
	}
}

func TestGetRandom(t *testing.T) {
	quotes := quote.Quotes{
		{Quote: "Quote 1", Author: "Author 1"},
		{Quote: "Quote 2", Author: "Author 2"},
	}

	quote, err := quotes.GetRandom()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if quote == nil || (quote.Quote != "Quote 1" && quote.Quote != "Quote 2") {
		t.Errorf("Expected a valid quote, got '%v'", quote)
	}
}

func TestGetRandom_Empty(t *testing.T) {
	var quotes quote.Quotes

	quote, err := quotes.GetRandom()

	if err != nil {
		t.Errorf("Expected no error for empty quotes, got %v", err)
	}

	if quote != nil {
		t.Errorf("Expected nil for empty quotes, got '%v'", quote)
	}
}

func TestQuotesJSON(t *testing.T) {
	quotes := quote.Quotes{
		{Quote: "Quote 1", Author: "Author 1"},
		{Quote: "Quote 2", Author: "Author 2"},
	}

	jsonStr, err := quotes.Json()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := `[{"quote":"Quote 1","author":"Author 1"},{"quote":"Quote 2","author":"Author 2"}]`

	if jsonStr != expected {
		t.Errorf("Expected '%s', got '%s'", expected, jsonStr)
	}
}

func TestQuoteJSON(t *testing.T) {
	quote := quote.Quote{
		Quote:  "Test Quote",
		Author: "Tester",
	}

	jsonStr, err := quote.Json()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := `{"quote":"Test Quote","author":"Tester"}`

	if jsonStr != expected {
		t.Errorf("Expected '%s', got '%s'", expected, jsonStr)
	}
}

func TestSave(t *testing.T) {
	quotes := quote.Quotes{
		{Quote: "Quote 1", Author: "Author 1"},
		{Quote: "Quote 2", Author: "Author"},
	}

	filename := "test_save_quotes.json"

	err := quotes.Save(filename)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	defer os.Remove(filename)

	var loadedQuotes quote.Quotes
	err = loadedQuotes.Load(filename)

	if err != nil {
		t.Errorf("Expected no error loading saved file, got %v", err)
	}

	if loadedQuotes.Count() != quotes.Count() {
		t.Errorf("Expected count %d, got %d", quotes.Count(), loadedQuotes.Count())
	}

	for i := range quotes {
		if loadedQuotes[i] != quotes[i] {
			t.Errorf("Mismatch at index %d: expected '%v', got '%v'", i, quotes[i], loadedQuotes[i])
		}
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	var quotes quote.Quotes
	err := quotes.Load("non_existent_file.json")

	if err == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	quotes := quote.Quotes{
		{Quote: "Quote 1", Author: "Author 1"},
	}

	err := quotes.Save("/invalid_path/test_save_quotes.json")

	if err == nil {
		t.Errorf("Expected error for invalid path, got nil")
	}
}

func TestAddQuote(t *testing.T) {
	var quotes quote.Quotes

	quotes.Add(quote.Quote{Quote: "New Quote", Author: "New Author"})

	if quotes.Count() != 1 {
		t.Errorf("Expected count 1, got %d", quotes.Count())
	}

	if quotes[0].Quote != "New Quote" || quotes[0].Author != "New Author" {
		t.Errorf("Quote data mismatch")
	}
}

func TestRemoveByID(t *testing.T) {
	quotes := quote.Quotes{
		{Quote: "Quote 1", Author: "Author 1"},
		{Quote: "Quote 2", Author: "Author 2"},
	}

	removed := quotes.RemoveByID(0)

	if !removed {
		t.Errorf("Expected true for successful removal, got false")
	}

	if quotes.Count() != 1 {
		t.Errorf("Expected count 1 after removal, got %d", quotes.Count())
	}

	if quotes[0].Quote != "Quote 2" {
		t.Errorf("Expected remaining quote to be 'Quote 2', got '%s'", quotes[0].Quote)
	}
}

func TestRemoveByID_Empty(t *testing.T) {
	var quotes quote.Quotes

	removed := quotes.RemoveByID(0)

	if removed {
		t.Errorf("Expected false for out-of-bounds from empty quotes, got true")
	}
}

func TestUpdateByID(t *testing.T) {
	quotes := quote.Quotes{
		{Quote: "Quote 1", Author: "Author 1"},
		{Quote: "Quote 2", Author: "Author 2"},
	}

	updated := quotes.UpdateByID(1, quote.Quote{Quote: "Updated Quote", Author: "Updated Author"})

	if !updated {
		t.Errorf("Expected true for successful update, got false")
	}

	if quotes[1].Quote != "Updated Quote" || quotes[1].Author != "Updated Author" {
		t.Errorf("Quote data mismatch after update")
	}
}
