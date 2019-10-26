package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

var (
	errRepository = errors.New("cannot get repository")
)

func newGithubClient(token string, timeout time.Duration) *github.Client {
	var client *http.Client
	if token == "" {
		client = &http.Client{}
	} else {
		client = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			),
		)
	}

	client.Timeout = timeout
	return github.NewClient(client)
}

func getRepository(token string) (string, error) {
	client := newGithubClient(token, time.Second)

	ctx := context.Background()
	opts := &github.ListOptions{}

	repos, _, err := client.Apps.ListRepos(ctx, opts)
	if err != nil {
		log.Printf("cannot get repositories: %s", err.Error())
		return "", errRepository
	}

	if len(repos) != 1 {
		log.Printf("incorrect number of repositories: %v", len(repos))
		return "", errRepository
	}

	if repos[0] == nil {
		log.Printf("empty repository")
		return "", errRepository
	}

	return repos[0].GetFullName(), nil
}

func main() {
	fmt.Println("codecov-action proxy")

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("No GITHUB_TOKEN env variable found")
	}

	repo, err := getRepository(token)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("repository: %v\n", repo)
}
