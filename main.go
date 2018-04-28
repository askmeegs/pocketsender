package main

import (
	"log"

	"io/ioutil"
	"os"

	"github.com/m-okeefe/pocketsender/pkg/pocketsender"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "Load configuration from file",
		},
	}

	app.Action = func(c *cli.Context) error {
		configPath := c.String("config")
		raw, err := ioutil.ReadFile(configPath)
		if err != nil {
			return err
		}

		cfg, err := pocketsender.NewConfig(raw)
		if err != nil {
			return err
		}

		ps, err := pocketsender.NewPocketSender(cfg)
		if err != nil {
			return err
		}

		return ps.Check()
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
