package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	message := "Call format: fileserve <port> <rootdir>"
	if len(os.Args) != 3 {
		fmt.Println(message)
		return
	}

	sport := os.Args[1]
	root := os.Args[2]

	_, err := strconv.Atoi(sport)
	if err != nil {
		fmt.Printf("Invalid port number: %v\n", sport)
		fmt.Println(message)
	}

	sport = ":" + sport

	_, err = os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Path %v does not exist\n", root)
		} else {
			fmt.Printf("Path %v not found: %v\n", root, err)
		}
	}

	fmt.Printf("Server directory %v served on port %v\n", root, sport)

	http.ListenAndServe(sport, http.FileServer(http.Dir(root)))

}
