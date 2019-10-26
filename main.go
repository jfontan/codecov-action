package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v28/github"
	"github.com/sanity-io/litter"
	"golang.org/x/oauth2"
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

func main() {
	fmt.Println("codecov-action proxy")

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("No GITHUB_TOKEN env variable found")
	}

	client := newGithubClient(token, time.Second)

	ctx := context.Background()
	opts := &github.ListOptions{}

	repos, response, err := client.Apps.ListRepos(ctx, opts)
	if err != nil {
		log.Fatal("Cannot get repositories", err)
	}

	fmt.Println("RESPONSE")
	litter.Dump(response)

	fmt.Println("REPOS")
	litter.Dump(repos)
}
