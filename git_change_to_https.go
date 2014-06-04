package main

import (
	"os"
	"path"
	"path/filepath"
)

func traverse(p string) {

}

func change_config(p string) {

}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please give a root path as an argument")
		return
	}

	traverse(os.Args[1])
}
