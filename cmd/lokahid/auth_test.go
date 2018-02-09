package main

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func TestAuth(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "everything is okay :)", http.StatusOK)
	})

	cases := []struct {
		name                   string
		serverUser, serverPass string
		clientUser, clientPass string
	}{}

	for _, cs := range cases {
		ts := httptest.NewServer(auth(cs.serverUser, cs.serverPass)(h))
		defer ts.Close()

		req := httptest.NewRequest("GET", "/", nil)

		req.Header.Add("Authorization", "Basic "+basicAuth(cs.clientUser, cs.clientPass))

		resp, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected resp.StatusCode to be %v, got: %v", http.StatusOK, resp.StatusCode)
		}
	}
}
