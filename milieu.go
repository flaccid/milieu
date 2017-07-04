package main

import (
	"fmt"
	"os"
	"time"

	"github.com/flaccid/milieu/walk"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "milieu"
	app.Version = "v0.0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Chris Fordham",
			Email: "chris@fordham-nagy.id.au",
		},
	}
	app.Copyright = "(c) 2016 Chris Fordham"
	app.Usage = "A tool to look at your source tree and advise on any lack of committment"
	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			fmt.Println("Usage: milieu [global options] <sourcetree_root>")
			os.Exit(1)
		}
		location := c.Args().Get(0)

		walk.Walk(location)

		return nil
	}

	app.Run(os.Args)
}
