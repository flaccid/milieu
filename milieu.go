package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
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
		fmt.Println("scanning from " + location)

		err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				_, err := os.Stat(path + "/.git/config")
				if err == nil {
					var (
						cmdOut []byte
					)
					cmdName := "git"
					cmdArgs := []string{"status"}
					os.Chdir(path)
					if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
						fmt.Fprintln(os.Stderr, err)
						return nil
					}
					if strings.Contains(string(cmdOut), "nothing to commit, working tree clean") {
						fmt.Println("> (", color.GreenString("ok"), ")", filepath.Base(path))
					} else {
						fmt.Println("> " + color.CyanString(filepath.Base(path)), "    ", path)
						fmt.Println(string(cmdOut))
						fmt.Print("Press 'Enter' to continue...")
						bufio.NewReader(os.Stdin).ReadBytes('\n')
					}
				}
				if os.IsNotExist(err) {
					return nil
				}
				return nil
			}
			return nil
		})
		if err != nil {
			fmt.Printf("walk error [%v]\n", err)
		}
		return nil
	}

	app.Run(os.Args)
}
