package main

import (
	"fmt"
	"github.com/flaccid/milieu"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"time"
)

func beforeApp(c *cli.Context) error {
	fmt.Printf("milieu %s\n", milieu.VERSION)

	switch c.GlobalString("log-level") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "":
		log.SetLevel(log.InfoLevel)
	default:
		log.Fatalf("%s is an invalid log level", c.GlobalString("log-level"))
	}

	log.Info("using log level " + log.GetLevel().String())

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "milieu"
	app.Version = milieu.VERSION
	app.Compiled = time.Now()
	app.Copyright = "(c) 2016 Chris Fordham"
	app.Usage = "A tool to look at your source tree and advise on any lack of commitment"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "prompt, p",
			Usage: "prompt",
		},
	}
	app.Run(os.Args)
}

func start(c *cli.Context) error {
	if c.NArg() < 1 {
		fmt.Println("Usage: milieu [global options] <sourcetree_root>")
		os.Exit(1)
	}
	location := c.Args().Get(0)

	milieu.Walk(location, c.Bool("prompt"))

	return nil
}
