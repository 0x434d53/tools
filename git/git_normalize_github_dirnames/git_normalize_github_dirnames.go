package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/0x434D53/tools/git/lib"
)

func dirExists(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func Rename(gi []lib.GitInfos) {
	for _, info := range gi {
		oldDir := path.Base(info.Path)
		newDir := info.Username + "_" + info.Projectname
		newDir = strings.TrimSuffix(newDir, ".git")
		fmt.Println(oldDir, newDir)

		if oldDir != newDir {

			basePath := path.Dir(info.Path)
			newPath := path.Join(basePath, newDir)

			if dirExists(newPath) {
				err := os.RemoveAll(newPath)

				if err != nil {
					fmt.Printf("Could not delete %v\n", newPath)
				}
			}

			err := os.Rename(info.Path, newPath)

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please give a root path as an argument")
		return
	}

	gitInfos, err := lib.CollectGitRepositories(os.Args[1])

	if err != nil {
		panic(err)
	}
	Rename(gitInfos)
}
