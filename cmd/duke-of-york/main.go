package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var should500 bool
var shouldLock sync.RWMutex

func main() {
	if p := os.Getenv("PORT"); p == "" {
		os.Setenv("PORT", "9001")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		shouldLock.RLock()
		defer shouldLock.RUnlock()

		if should500 {
			http.Error(w, "HTTP 500: Something really bad happened", http.StatusInternalServerError)
		} else {
			http.Error(w, "HTTP 200: Everything is OK", http.StatusOK)
		}
	})

	go func() {
		t := time.NewTicker(time.Minute)
		for {
			select {
			case <-t.C:
				shouldLock.Lock()
				should500 = !should500

				log.Printf("will return 500: %v", should500)

				shouldLock.Unlock()
			}
		}
	}()

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
