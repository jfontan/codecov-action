package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/v28/github"
	"github.com/sanity-io/litter"
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

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("new connection: %v", r.URL)
	r.URL.Query().Set("token", "LOLO")
	log.Printf("rewritten connection: %v", r.URL)

	request, err := http.NewRequest("POST", "https://codecov.io/upload/v4", nil)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	params := r.URL.Query()
	params.Set("token", "f4d08eb3-7b27-432b-955b-2c1b71d2b800")
	request.URL.RawQuery = params.Encode()

	println("query", request.URL.RawQuery)

	log.Printf("rewritten connection: %v", request.URL)
	litter.Dump(request.Form)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("error connecting: %v", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	println("body", string(body))

	if resp.StatusCode != 200 {
		log.Fatalf("error connecting: %v", resp.StatusCode)
		http.Error(w, http.StatusText(resp.StatusCode), resp.StatusCode)
		return
	}

	// _, err = io.Copy(w, resp.Body)
	_, err = w.Write(body)
	if err != nil {
		log.Fatalf("error writting response: %v", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func main() {
	fmt.Println("codecov-action proxy")

	// token := os.Getenv("GITHUB_TOKEN")
	// if token == "" {
	// 	log.Fatal("no GITHUB_TOKEN env variable found")
	// }

	// repo, err := getRepository(token)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("repository: %v\n", repo)

	http.HandleFunc("/upload/v4", handler)
	log.Fatal(http.ListenAndServe(":8808", nil))
}
