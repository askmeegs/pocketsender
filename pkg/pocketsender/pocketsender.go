package pocketsender

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	log.Println("Instantiating Pocket API client...")
	ps.PocketClient = pocketapi.NewClient(ps.PocketConsumerKey, ps.PocketAccessToken)

	log.Println("Retrieving unread Pocket articles for account...")
	retrieval, err := ps.PocketClient.Retrieve(&pocketapi.RetrieveOption{State: pocketapi.StateUnread})
	if err != nil {
		return err
	}

	log.Printf("Got %d unread pocket articles, emailing...\n", len(retrieval.List))
	i := 1
	for _, item := range retrieval.List {
		log.Printf("\n\n ---> Article %d/%d:\n", i, len(retrieval.List))
		err := ps.emailAndArchive(item)
		if err != nil {
			return err
		}
		i++
	}
	return nil
}

func (ps *PocketSender) emailAndArchive(item pocketapi.Item) error {

	url := strings.Replace(item.ResolvedURL, "https://", "http://", 1)

	// Get the html at that url --> convert to pdf. Save PDF locally.
	log.Printf("Converting to pdf: %s\n", url)
	pdfPath, err := ps.savePDF(url, item.ResolvedTitle)
	if err != nil {
		return err
	}

	// Send an email to kindle with this PDF attached
	err = ps.emailKindle(pdfPath)
	if err != nil {
		return err
	}

	// Remove temporary pdf from local fs
	err = os.Remove(pdfPath)
	if err != nil {
		return err
	}

	// Tell pocket to archive this article so that we don't send it again
	return ps.archiveInPocket(item.ItemID)
}

func (ps *PocketSender) savePDF(url string, title string) (string, error) {
	pdfg, err := htmltopdf.NewPDFGenerator()
	if err != nil {
		return "", fmt.Errorf("Could not instantiate pdf generator: %v", err)
	}

	pdfg.AddPage(htmltopdf.NewPage(url))

	err = pdfg.Create()
	if err != nil {
		log.Printf("[warning] PDF generator returned an error: %v\n", err)
	}

	pdfName := ps.generatePdfName(title)

	saveAt := fmt.Sprintf("./pdf/%s.pdf", pdfName)
	err = pdfg.WriteFile(saveAt)
	if err != nil {
		return "", fmt.Errorf("Could not write pdf to file: %v", err)
	}
	return saveAt, nil
}

func (ps *PocketSender) emailKindle(pdfPath string) error {
	log.Println("Emailing PDF to kindle...")
	m := gomail.NewMessage()
	m.SetHeader("From", ps.FromEmail)
	m.SetHeader("To", ps.KindleEmail)
	m.SetHeader("Subject", "convert")
	m.SetBody("text/html", "<p><3 from pocketsender</p>")
	m.Attach(pdfPath)

	googleUsername := strings.Split(ps.FromEmail, "@")[0]
	d := gomail.NewDialer("smtp.gmail.com", 587, googleUsername, ps.FromEmailPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	log.Println("Sent!")
	return nil
}

func (ps *PocketSender) generatePdfName(title string) string {
	words := strings.Split(title, " ")
	if len(words) < 5 {
		return strings.Join(words, "-")
	}
	return strings.Join(words[:5], "-")
}

func (ps *PocketSender) archiveInPocket(id int) error {
	archiveAction := pocketapi.NewArchiveAction(id)
	_, err := ps.PocketClient.Modify(archiveAction)
	if err != nil {
		return err
	}
	return nil
}
