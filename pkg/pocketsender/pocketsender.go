package pocketsender

import (
	"encoding/json"
	"fmt"
	"strings"

	htmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	pocketapi "github.com/motemen/go-pocket/api"
	"gopkg.in/gomail.v2"
)

type Config struct {
	PocketUsername    string
	PocketAccessToken string
	PocketConsumerKey string
	FromEmail         string
	FromEmailPassword string
	KindleEmail       string
}

type PocketSender struct {
	PocketUsername    string
	PocketAccessToken string
	PocketConsumerKey string
	FromEmail         string
	FromEmailPassword string
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
		FromEmailPassword: cfg.FromEmailPassword,
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
		if item.ItemID == 2128743607 {
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

	// url := strings.Replace(item.ResolvedURL, "https://", "http://", 1)
	//
	// // Get the html at that url --> convert to pdf. Save PDF locally.
	// fmt.Printf("\n\nConverting to pdf: %s\n", url)
	// pdfPath, err := ps.savePDF(url, item.ItemID)
	// if err != nil {
	// 	return err
	// }

	// Send an email to kindle where at pdfPath
	pdfPath := "./pdf/pocketsender-2128743607.pdf"
	err := ps.emailKindle(pdfPath)
	if err != nil {
		return err
	}
	return nil

	// Scrub pdf from filesystem
	// err = os.Remove(pdfPath)
	// if err != nil {
	// 	return err
	// }
	//
	// // Tell pocket to archive this article so that we don't send it again
	// return ps.archiveInPocket(item.ItemID)
}

func (ps *PocketSender) savePDF(url string, id int) (string, error) {
	fmt.Println("Instantiating pdf generator...")
	pdfg, err := htmltopdf.NewPDFGenerator()
	if err != nil {
		return "", fmt.Errorf("Could not instantiate pdf generator: %v", err)
	}

	fmt.Println("Adding page to generator...")
	pdfg.AddPage(htmltopdf.NewPage(url))

	fmt.Println("Creating PDF...")
	err = pdfg.Create()
	if err != nil {
		fmt.Printf("Got err but ignoring - Could not create pdf: %v\n", err)
	}
	saveAt := fmt.Sprintf("./pdf/pocketsender-%d.pdf", id)
	err = pdfg.WriteFile(saveAt)
	if err != nil {
		return "", fmt.Errorf("Could not write pdf to file: %v", err)
	}
	return saveAt, nil
}

func (ps *PocketSender) emailKindle(pdfPath string) error {
	fmt.Println("Composing email message for this pdf...")
	m := gomail.NewMessage()
	m.SetHeader("From", ps.FromEmail)
	m.SetHeader("To", ps.KindleEmail)
	m.SetHeader("Subject", "convert")
	m.Attach(pdfPath)

	fmt.Println("Dialing smtp...")

	googleUsername := strings.Split(ps.FromEmail, "@")[0]
	d := gomail.NewDialer("smtp.gmail.com", 587, googleUsername, ps.FromEmailPassword)

	fmt.Println("Sending email...")
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	fmt.Println("Sent!")
	return nil
}

func (ps *PocketSender) archiveInPocket(id int) error {
	return nil
}
