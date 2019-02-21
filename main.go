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
			Name: "organization",
		},
	}
	app.Action = func(c *cli.Context) error {
		baseUrl := c.String("github-base-url")
		uploadUrl := c.String("github-upload-url")
		org := c.String("organization")

		if os.Getenv("GITHUB_BASE_URL") != "" {
			baseUrl = os.Getenv("GITHUB_BASE_URL")
		}

		if os.Getenv("GITHUB_UPLOAD_URL") != "" {
			uploadUrl = os.Getenv("GITHUB_UPLOAD_URL")
		}

		if os.Getenv("ORGANIZATION") != "" {
			org = os.Getenv("ORGANIZATION")
		}

		return server.Start(baseUrl, uploadUrl, org)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
