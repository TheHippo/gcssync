package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gcssync"
	app.Usage = "Sync files with Google Cloud Storage"
	app.Commands = []cli.Command{
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "List remote files",
			Action:    listFiles,
		},
	}
	app.Run(os.Args)
}

func listFiles(c *cli.Context) {
	fmt.Println("List files")
}
