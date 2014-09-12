package main

import (
	"fmt"
	"github.com/0x434D53/cms_go/filecrypt"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Call Format: decryptFile <password> <oldfilename> <newfilename>")
		return
	}

	password := os.Args[1]
	oldFilename := os.Args[2]
	newFilename := os.Args[3]

	key, err := filecrypt.CreateKeyFromPassword(password)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = filecrypt.DecryptFile(oldFilename, newFilename, key)

	if err != nil {
		fmt.Println(err)
	}
}
