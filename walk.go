package walk

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/tj/go-spin"
)

func Round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}

// Walk - step through a directory recursively looking for gits
func Walk(location string, prompt bool) {
	fmt.Println("scanning " + location)

	files := 0
	total := 0
	clean := 0
	modified := 0

	s := spin.New()
	s.Set(spin.Box3)

	start := time.Now()

	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		files++
		if info.IsDir() {
			_, err := os.Stat(path + "/.git/config")
			if err == nil {
				total++
				fmt.Printf("\r\033[36m\033[m %s ", s.Next())
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
					clean++
					fmt.Println("(" + color.GreenString("ok") + ") " + filepath.Base(path))
				} else {
					modified++
					fmt.Println("(" + color.YellowString("!m") + ") " + color.CyanString(filepath.Base(path)) + "    " + path)
					fmt.Println(string(cmdOut))
					if prompt {
						fmt.Print("Press 'Enter' to continue...")
						bufio.NewReader(os.Stdin).ReadBytes('\n')
					}
				}
			}
			if os.IsNotExist(err) {
				return nil
			}
			return nil
		} else {
			fmt.Printf("\r\033[36m\033[m %s ", s.Next())
		}
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	elapsed := time.Since(start)

	fmt.Println(color.GreenString("finished"))
	fmt.Println("\nSummary")
	fmt.Println(fmt.Sprintf("\ttotal: %v", total))
	fmt.Println(fmt.Sprintf("\tclean: %v", clean))
	fmt.Println(fmt.Sprintf("\tmodified: %v", modified))
	fmt.Println(fmt.Sprintf("\tscore: %v%%", Round(float64(clean)/float64(total)*100, 0.05)))
	fmt.Println(fmt.Sprintf("\tfiles walked: %v", files))
	fmt.Println(fmt.Sprintf("\tscan time: %v", elapsed))
}
