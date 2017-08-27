package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubURLType int

const (
	ssh GithubURLType = iota
	git
	https
)

func main() {
	accessToken := os.Getenv("GithubAccessToken")

	if accessToken == "" {
		log.Fatal("Not an Access Token found as GithubAccessToken in Environment")
	}

	goget := flag.Bool("goget", false, "should generate a git clone or go get -u for the reps")
	flag.Parse()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	repos, err := getStarredReposForUser(client, https)

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range repos {
		if *goget {
			fmt.Printf("go get %s/...\n", r.geturl)
		} else {
			fmt.Printf("git clone %s %s_%s\n", r.sshurl, r.user, r.repo)
		}
	}
}

type repo struct {
	sshurl string
	geturl string
	user   string
	repo   string
}

func getStarredReposForUserPage(client *github.Client, urlType GithubURLType, page int) ([]repo, int, error) {
	options := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{Page: page, PerPage: 50},
	}

	starredRepos, response, err := client.Activity.ListStarred(context.Background(), "", options)

	if err != nil {
		return nil, 0, err
	}

	repos := make([]repo, 0)

	for _, r := range starredRepos {
		newr := repo{sshurl: *r.Repository.SSHURL, geturl: fmt.Sprintf("github.com/%s", *r.Repository.FullName), user: *r.Repository.Owner.Login, repo: *r.Repository.Name}

		repos = append(repos, newr)
	}

	return repos, response.LastPage, nil
}

func getStarredReposForUser(client *github.Client, urlType GithubURLType) ([]repo, error) {
	repos, lastPage, err := getStarredReposForUserPage(client, https, 0)

	if err != nil {
		return nil, err
	}

	if lastPage <= 0 {
		return repos, nil
	}

	for i := 1; i <= lastPage; i++ {
		reposPage, _, err := getStarredReposForUserPage(client, https, i)

		if err != nil {
			return nil, err
		}

		repos = append(repos, reposPage...)
	}

	return repos, nil
}
