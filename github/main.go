package main

import (
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "64b91e2e931f33ff53b4f20201cd02f399a5fbfe"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	repos, _, err := client.Repositories.List("", nil)

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range repos {
		fmt.Println(*r.Name)
	}
}
