package pocketsender

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Config struct {
	PocketEmail string
	FromEmail   string
	KindleEmail string
}

type PocketSender struct {
	PocketEmail string
	FromEmail   string
	KindleEmail string
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
		PocketEmail: cfg.PocketEmail,
		FromEmail:   cfg.FromEmail,
		KindleEmail: cfg.KindleEmail,
	}, nil
}

func (ps *PocketSender) Watch() error {
	log.Printf("%#v", ps)
	time.Sleep(time.Hour)
	return nil
}
