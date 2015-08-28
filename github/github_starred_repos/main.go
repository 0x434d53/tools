package main

import (
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

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	repos, err := getStarredReposForUser(client, https)

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range repos {
		fmt.Printf("git clone %v\n", r)
	}
}

func getStarredReposForUserPage(client *github.Client, urlType GithubURLType, page int) ([]string, int, error) {
	options := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{Page: page, PerPage: 50},
	}

	starredRepos, response, err := client.Activity.ListStarred("", options)

	if err != nil {
		return nil, 0, err
	}

	repos := make([]string, 0)

	for _, r := range starredRepos {
		repos = append(repos, *r.Repository.SSHURL)
	}

	return repos, response.LastPage, nil
}

func getStarredReposForUser(client *github.Client, urlType GithubURLType) ([]string, error) {
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
