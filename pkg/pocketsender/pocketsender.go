package pocketsender

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	pocketapi "github.com/motemen/go-pocket/api"
	"github.com/motemen/go-pocket/auth"
)

type Config struct {
	PocketEmail       string
	PocketConsumerKey string
	FromEmail         string
	KindleEmail       string
}

type PocketSender struct {
	PocketEmail       string
	PocketConsumerKey string
	FromEmail         string
	KindleEmail       string
}

func NewConfig(input []byte) (*Config, error) {
	var cfg Config

	err := json.Unmarshal(input, &cfg)
	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}

func NewPocketSender(cfg *Config) (*PocketSender, error) {
	if len(cfg.PocketEmail) == 0 || len(cfg.FromEmail) == 0 || len(cfg.KindleEmail) == 0 {
		return &PocketSender{}, fmt.Errorf("Must supply pocket, from, and kindle emails")
	}
	return &PocketSender{
		PocketEmail:       cfg.PocketEmail,
		PocketConsumerKey: cfg.PocketConsumerKey,
		FromEmail:         cfg.FromEmail,
		KindleEmail:       cfg.KindleEmail,
	}, nil
}

func (ps *PocketSender) Check() error {
	pocketAccess, err := obtainAccessToken(ps.PocketConsumerKey)
	if err != nil {
		return err
	}
	pocketClient := pocketapi.NewClient(ps.PocketConsumerKey, pocketAccess.AccessToken)

	retrieval, err := pocketClient.Retrieve(&pocketapi.RetrieveOption{})
	if err != nil {
		return err
	}

	pp, _ := json.MarshalIndent(retrieval, "", "  ")
	fmt.Println(pp)

	return nil
}

func obtainAccessToken(consumerKey string) (*auth.Authorization, error) {
	ch := make(chan struct{})
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/favicon.ico" {
				http.Error(w, "Not Found", 404)
				return
			}

			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(w, "Authorized.")
			ch <- struct{}{}
		}))
	defer ts.Close()

	redirectURL := ts.URL

	requestToken, err := auth.ObtainRequestToken(consumerKey, redirectURL)
	if err != nil {
		return nil, err
	}

	url := auth.GenerateAuthorizationURL(requestToken, redirectURL)
	fmt.Println(url)

	<-ch

	return auth.ObtainAccessToken(consumerKey, requestToken)
}
