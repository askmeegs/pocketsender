package main

import (
	"fmt"
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

		fmt.Println("****************************************************************************************\nP O C K E T S E N D E R             v0.0.1\n****************************************************************************************")

		if _, err := os.Stat("./pdf/"); os.IsNotExist(err) {
			err := os.Mkdir("./pdf/", 0777)
			if err != nil {
				return err
			}
		}

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
