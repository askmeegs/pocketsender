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
	for _, item := range retrieval.List {
		if item.ItemID == 2011165197 {
			err := ps.emailAndArchive(item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ps *PocketSender) emailAndArchive(item pocketapi.Item) error {
	// Get the html at that url --> convert to pdf. Save PDF locally.
	log.Printf("\n\n\nConverting to pdf: %s\n", item.GivenURL)
	pdfPath, err := ps.savePDF(item.GivenURL, item.ResolvedTitle)
	if err != nil {
		return err
	}

	// Send an email to kindle where at pdfPath
	err = ps.emailKindle(pdfPath)
	if err != nil {
		return err
	}
	return nil

	// Remove pdf from local fs
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
	//
	// resolvedUrl, err := ps.resolveUrl(url) // handle any 301 / redirects
	// if err != nil {
	// 	return "", err
	// }
	// fmt.Printf("%s ---> resolved to ---> %s\n", url, resolvedUrl)

	pdfg.AddPage(htmltopdf.NewPage(url))

	err = pdfg.Create()
	if err != nil {
		return "", fmt.Errorf("Could not create pdf: %v\n", err)
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
	result, err := ps.PocketClient.Modify(archiveAction)
	if err != nil {
		return err
	}
	log.Printf("Archived item [%d] in Pocket. StatusCode=%d\n", id, result.Status)
	return nil
}

//
// func (ps *PocketSender) resolveUrl(url string) (string, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return url, err
// 	}
// 	defer resp.Body.Close()
//
// 	fmt.Printf("%#v", resp)
// 	pp, _ := json.MarshalIndent(&resp, "", "  ")
// 	log.Printf("\n %s\n\n", string(pp))
//
// 	_, err = resp.Location()
// 	if err != nil {
// 		return url, err
// 	}
// 	return url, nil
// 	// return location.String(), nil
// }
