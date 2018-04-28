package pocketsender

import (
	"encoding/json"
	"fmt"

	pocketapi "github.com/motemen/go-pocket/api"
)

type Config struct {
	PocketUsername    string
	PocketAccessToken string
	PocketConsumerKey string
	FromEmail         string
	KindleEmail       string
}

type PocketSender struct {
	PocketUsername    string
	PocketAccessToken string
	PocketConsumerKey string
	FromEmail         string
	KindleEmail       string

	PocketClient *pocketapi.Client
	sent         map[string]bool
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
	return &PocketSender{
		PocketUsername:    cfg.PocketUsername,
		PocketAccessToken: cfg.PocketAccessToken,
		PocketConsumerKey: cfg.PocketConsumerKey,
		FromEmail:         cfg.FromEmail,
		KindleEmail:       cfg.KindleEmail,
	}, nil
}

func (ps *PocketSender) Check() error {

	fmt.Println("Instantiating Pocket API client...")
	ps.PocketClient = pocketapi.NewClient(ps.PocketConsumerKey, ps.PocketAccessToken)

	fmt.Println("Retrieving unread Pocket articles for account...")
	retrieval, err := ps.PocketClient.Retrieve(&pocketapi.RetrieveOption{State: pocketapi.StateUnread})
	if err != nil {
		return err
	}

	fmt.Printf("Got %d unread pocket articles, emailing...\n", len(retrieval.List))
	for _, item := range retrieval.List {
		if item.ItemID == 1990618385 {
			err := ps.emailAndArchive(item)
			if err != nil {
				return err
			}
		}
	}

	// fmt.Printf("Successfully emailed %d articles.\n", len(retrieval.List))
	return nil
}

func (ps *PocketSender) emailAndArchive(item pocketapi.Item) error {
	// Extract url
	url := item.ResolvedURL

	err := ps.SavePDF(url, item.ItemID)
	if err != nil {
		return err
	}

	// Get the html at that url --> convert to pdf. Save PDF locally.

	// Send an email to kindle where this PDF is the attachment

	// Scrub pdf from filesystem

	// Tell pocket to archive this article so that we don't send it again

	return nil
}

func (ps *PocketSender) SavePDF(url string, id int) error {
	return nil
}
