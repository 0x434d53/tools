package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatal("Access token not found in environment")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(accessToken)})

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opts := &github.RepositoryListByAuthenticatedUserOptions{}

	repos, _, err := client.Repositories.ListByAuthenticatedUser(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	reposPerOwner := map[string][]string{}

	for _, r := range repos {
		if _, ok := reposPerOwner[*r.Owner.Login]; !ok {
			reposPerOwner[*r.Owner.Login] = []string{}
		}
		reposPerOwner[*r.Owner.Login] = append(reposPerOwner[*r.Owner.Login], *r.SSHURL)
	}

	for dir, rs := range reposPerOwner {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			_ = os.Mkdir(dir, os.ModePerm)
		}

		if err := os.Chdir(dir); err != nil {
			log.Fatal(err)
		}

		for _, r := range rs {
			u := strings.TrimSuffix(r, ".git")
			fmt.Printf("cloning %s\n", u)
			cmd := exec.Command("git", "clone", u)
			err := cmd.Run()
			if err != nil {
				log.Print(err)
			}
		}

		if err := os.Chdir(".."); err != nil {
			log.Fatal(err)
		}
	}
}
