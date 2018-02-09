package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
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
		wantStatus             int
	}{
		{
			name:       "ok",
			serverUser: "shachi",
			serverPass: "orca",
			clientUser: "shachi",
			clientPass: "orca",
			wantStatus: http.StatusOK,
		},
		{
			name:       "not ok",
			serverUser: "shachi",
			serverPass: "orca",
			clientUser: "shachi",
			clientPass: "orcA",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			ts := httptest.NewServer(auth(cs.serverUser, cs.serverPass)(h))
			defer ts.Close()

			req := httptest.NewRequest("GET", ts.URL, nil)
			req.RequestURI = ""
			req.Host = ""
			req.Header.Add("Authorization", "Basic "+basicAuth(cs.clientUser, cs.clientPass))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != cs.wantStatus {
				data, err := httputil.DumpResponse(resp, true)
				if err != nil {
					t.Fatal(err)
				}

				fmt.Println(string(data))

				t.Fatalf("expected resp.StatusCode to be %d, got: %d", cs.wantStatus, resp.StatusCode)
			}
		})
	}
}
