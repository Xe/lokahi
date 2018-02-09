package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/caarlos0/env"
)

type dWebhook struct {
	Content  string `json:"content,omitifempty"`
	Username string `json:"username"`
}

type impl struct {
	cfg config
}

func (i impl) Handle(ctx context.Context, st *lokahi.CheckStatus) (*lokahi.Nil, error) {
	sendWebhook(i.cfg.WebhookURL, dWebhook{
		Username: "Lokahi",
		Content: fmt.Sprintf(
			"Service at %s is %s (%v in %v), playbook: <%s>",
			st.Check.Url,
			st.Check.State.String(),
			st.RespStatusCode,
			time.Duration(st.LastResponseTimeNanoseconds),
			st.Check.PlaybookUrl,
		),
	})

	return &lokahi.Nil{}, nil
}

type config struct {
	Port       string `env:"PORT" envDefault:"9001"`
	WebhookURL string `env:"WEBHOOK_URL,required"`
}

func sendWebhook(whurl string, dw dWebhook) error {
	data, err := json.Marshal(&dw)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(whurl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode/100 != 2 {
		io.Copy(os.Stderr, resp.Body)
		resp.Body.Close()
		return fmt.Errorf("status code was %v", resp.StatusCode)
	}

	return nil
}

func main() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	hs := lokahi.NewWebhookServer(impl{cfg: cfg}, nil)
	mux := http.NewServeMux()
	mux.Handle(lokahi.WebhookPathPrefix, hs)

	log.Printf("listening on http://0.0.0.0:%v", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
