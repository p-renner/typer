package quote

import (
	"fmt"
	"net/http"
	"strconv"
)

var quotes Quotes

func Server() {
	err := quotes.Load("quotes.json")

	if err != nil {
		fmt.Println("Error loading Quotes", err)
		return
	}

	router := http.NewServeMux()

	router.Handle("/quotes", applyMiddlewares(http.HandlerFunc(quotesHandler), loggerMiddleware, jsonMiddleware))
	router.Handle("/quote", applyMiddlewares(http.HandlerFunc(quoteHandler), loggerMiddleware, jsonMiddleware))
	router.Handle("/randomquote", applyMiddlewares(http.HandlerFunc(randomQuoteHandler), loggerMiddleware, jsonMiddleware))

	fmt.Println("Server started at :8080")
	err = http.ListenAndServe(":8080", router)

	if err != nil {
		fmt.Println("Error starting server", err)
		return
	}
}

func applyMiddlewares(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, r.Proto, r.UserAgent())
		next.ServeHTTP(w, r)
	})
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]

	if !ok || len(ids[0]) < 1 {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(ids[0])

	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	quote, err := quotes.GetByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	q, err := quote.Json()

	if err != nil {
		http.Error(w, "Error parsing JSON"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(q))
}

func randomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	quote, err := quotes.GetRandom()

	if err != nil {
		http.Error(w, "Error retrieving quote"+err.Error(), http.StatusInternalServerError)
		return
	}

	if quote == nil {
		http.Error(w, "No quotes available", http.StatusNotFound)
		return
	}

	q, err := quote.Json()

	if err != nil {
		http.Error(w, "Error parsing JSON"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(q))
}

func quotesHandler(w http.ResponseWriter, r *http.Request) {
	q, err := quotes.Json()

	if err != nil {
		http.Error(w, "Error parsing JSON"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(q))
}
