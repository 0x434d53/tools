package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	// var lastip string

	u := "http://checkip.dyndns.org"

	resp, err := http.Get(u)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	defer resp.Body.Close()

	c, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	r := regexp.MustCompile("[0-9.]+")
	res := r.Find(c)

	fmt.Println(string(res))
}
