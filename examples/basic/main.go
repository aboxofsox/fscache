package main

import (
	"log"
	"net/http"

	"github.com/aboxofsox/fscache"
)

func setValue(cache *fscache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		key := query.Get("key")
		value := query.Get("value")

		if key == "" || value == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		cache.Set(key, value)

		w.Write([]byte("OK"))
	}
}

func getValue(cache *fscache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		key := query.Get("key")

		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		value, ok := cache.Get(key)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return
		}

		w.Write([]byte(value.(string)))
	}
}

func main() {
	cache := fscache.NewCache("foo.gob")

	http.HandleFunc("/set", setValue(cache))
	http.HandleFunc("/get", getValue(cache))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
