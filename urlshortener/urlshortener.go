package urlshortener

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type ResponseMessage[T any] struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    T      `json:"data,omitempty"`
}

type ShortenerCreate struct {
	URL string `json:"url"`
}

func SecurityHeaders(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Content-Security-Policy", "frame-ancestors 'none'")
		w.Header().Add("X-Content-Type", "nosniff")
		w.Header().Add("X-Frame-Options", "DENY")

		f.ServeHTTP(w, r)
	}
}

func JsonContent(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		contentType := r.Header.Get("Content-Type")
		if contentType != "" && !strings.Contains(contentType, "json") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			json.NewEncoder(w).Encode(ResponseMessage[any]{
				Error: "Only supports `application/json`",
			})

			return
		}

		f.ServeHTTP(w, r)
	}
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func ChainMiddleware(f http.HandlerFunc, ms ...Middleware) http.HandlerFunc {
	ret := f
	for _, m := range ms {
		ret = m(ret)
	}

	return ret
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")

	s := Shortener{Hash: hash}
	if err := s.Get(); err != nil {
		var code int
		errMsg := err.Error()
		switch err {
		case ErrInvalidHash:
			code = http.StatusBadRequest
		case sql.ErrNoRows:
			code = http.StatusNotFound
			errMsg = "Path not mapped to a URL"
		default:
			code = http.StatusInternalServerError
			errMsg = "Internal Server Error"
		}

		w.WriteHeader(code)
		json.NewEncoder(w).Encode(ResponseMessage[any]{
			Error: errMsg,
		})

		return
	}

	http.Redirect(w, r, s.URL, http.StatusFound)
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(w)

	var createRequest ShortenerCreate
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e.Encode(ResponseMessage[any]{
			Error: "Cannot parse Shortener Creation request",
		})

		return
	}

	s := Shortener{URL: createRequest.URL}
	if err := s.Create(12); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e.Encode(ResponseMessage[any]{
			Error: "Internal server error",
		})

		return
	}

	w.WriteHeader(http.StatusCreated)
	e.Encode(ResponseMessage[Shortener]{
		Message: "Shortener creation successful",
		Data:    s,
	})
}
