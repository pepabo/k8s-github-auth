package main

import (
	"github.com/takaishi/k8s-github-auth/server"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "github-base-url",
		},
		cli.StringFlag{
			Name: "github-upload-url",
		},
		cli.StringFlag{
			Name: "team",
		},
	}
	app.Action = func(c *cli.Context) error {
		return server.Start(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
