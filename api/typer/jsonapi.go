package main

import (
	"net/http"
	"strconv"

	"typer/quote"
)

var quotes quote.Quotes

func main() {
	if err := quotes.Load("quotes.json"); err != nil {
		panic(err)
	}

	http.HandleFunc("/quote", quoteHandler)
	http.HandleFunc("/quotes", quotesHandler)

	http.ListenAndServe(":8080", nil)
}

// quoteHandler handles requests to /quote endpoint.
func quoteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getQuoteHandler(w, r)
	case http.MethodDelete:
		deleteQuoteHandler(w, r)
	case http.MethodPost:
		createQuoteHandler(w, r)
	case http.MethodPut:
		updateQuoteHandler(w, r)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// getQuoteHandler handles GET requests to /quote endpoint.
// It returns a random quote or a quote by ID in JSON format.
func getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("id")

	var quote *quote.Quote
	var err error

	if ids == "" {
		quote, err = quotes.GetRandom()
	} else {
		id, err := strconv.Atoi(ids)

		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		quote, err = quotes.GetByID(id)
	}

	if quote == nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Failed to get quote", http.StatusInternalServerError)
		return
	}

	json, err := quote.Json()

	if err != nil {
		http.Error(w, "Failed to serialize quote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json))
}

// deleteQuoteHandler handles DELETE requests to /quote endpoint.
// It deletes a quote by ID.
func deleteQuoteHandler(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("id")

	if ids == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(ids)

	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if deleted := quotes.RemoveByID(id); !deleted {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// createQuoteHandler handles POST requests to /quote endpoint.
// It creates a new quote.
func createQuoteHandler(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	author := r.URL.Query().Get("author")

	if text == "" || author == "" {
		http.Error(w, "Missing text or author", http.StatusBadRequest)
		return
	}

	quotes.Add(quote.Quote{
		Quote:  text,
		Author: author,
	})

	w.WriteHeader(http.StatusCreated)
}

// updateQuoteHandler handles PUT requests to /quote endpoint.
// It updates a quote by ID.
func updateQuoteHandler(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("id")
	text := r.URL.Query().Get("text")
	author := r.URL.Query().Get("author")

	if ids == "" || text == "" || author == "" {
		http.Error(w, "Missing id, text or author", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(ids)

	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if updated := quotes.UpdateByID(id, quote.Quote{
		Quote:  text,
		Author: author,
	}); !updated {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// quotesHandler handles requests to /quotes endpoint.
// It returns all quotes in JSON format.
func quotesHandler(w http.ResponseWriter, r *http.Request) {
	json, err := quotes.Json()

	if err != nil {
		http.Error(w, "Failed to serialize quotes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json))
}
