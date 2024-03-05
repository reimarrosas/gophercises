package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gophercises/urlshortener"
	"log"
	"net/http"
)

func main() {
	var port string
	flag.StringVar(&port, "port", ":8080", "Port of the Web App")
	flag.Parse()

	urlshortener.InitLogger("")
    if err := urlshortener.InitDB("shortener.db"); err != nil {
        panic(err)
    }

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{hash}", urlshortener.HandleRedirect)
	mux.HandleFunc("POST /{$}", urlshortener.HandlePost)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(urlshortener.ResponseMessage[any]{
			Error: fmt.Sprintf("%s not found", r.URL.Path),
		})
	})

	handler := urlshortener.ChainMiddleware(
		mux.ServeHTTP,
		urlshortener.SecurityHeaders,
		urlshortener.JsonContent,
        urlshortener.LoggerMiddleware,
	)

	log.Fatal(http.ListenAndServe(port, handler))
}
