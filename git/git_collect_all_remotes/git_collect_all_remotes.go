package main

import (
	"fmt"
	"github.com/0x434D53/tools/git/lib"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please give a root path as an argument")
		return
	}

	gitInfos, err := lib.CollectGitRepositories(os.Args[1])

	if err != nil {
		panic(err)
	}

	for _, gi := range gitInfos {
		fmt.Printf("git clone git@github.com:%v/%v\n", gi.Username, gi.Projectname)
	}
}
