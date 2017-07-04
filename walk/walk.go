package walk

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/tj/go-spin"
)

// Walk - step through a directory recursively looking for gits
func Walk(location string) {
	fmt.Println("scanning " + location)

	s := spin.New()
	s.Set(spin.Box3)

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
					fmt.Println("[", color.GreenString("ok"), "]", filepath.Base(path))
				} else {
					fmt.Println(color.CyanString(filepath.Base(path)), "    ", path)
					fmt.Println(string(cmdOut))
					fmt.Print("Press 'Enter' to continue...")
					bufio.NewReader(os.Stdin).ReadBytes('\n')
				}
			}
			if os.IsNotExist(err) {
				return nil
			}
			return nil
		} else {
			fmt.Printf("\r  \033[36m\033[m %s ", s.Next())
		}
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}

	fmt.Println("[", color.GreenString("finished"), "]")
}
