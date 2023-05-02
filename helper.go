package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

type PullRequest struct {
	Title      string `json:"title"`
	URL        string `json:"html_url"`
	Created_At string `json:"created_at"`
	User       struct {
		Login string `json:"login"`
	}
}

// Make a request to the GitHub API to retrieve the list of pull requests for a particular repo and particular page.
// Send the result back as an array of PullRequest
func SendHttpRequest(repo string, page int, authToken string) ([]PullRequest, error) {
	url := "https://api.github.com/repos/" + GetRepoOwner() + "/" + repo + "/pulls?state=open&page=" + strconv.Itoa(page)

	fmt.Println("Sending request for " + url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []PullRequest{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []PullRequest{}, err
	}

	// Parse the response body and extract the relevant information
	defer resp.Body.Close()
	var pullRequests []PullRequest
	err = json.NewDecoder(resp.Body).Decode(&pullRequests)
	if err != nil {
		return []PullRequest{}, err
	}

	return pullRequests, nil
}

func getPullRequests(repo string, authToken string) []PullRequest {
	// Make requests to the GitHub API to retrieve the list of pull requests for all repos and range of pages.
	var combinedPullRequests []PullRequest

	teammates := GetTeammates()

	for pageNumber := 1; pageNumber <= GetNumberOfPages(); pageNumber++ {
		time.Sleep(1 * time.Second) // sleep for 1 sec so that we don't send requests too many at once.
		prs, err := SendHttpRequest(repo, pageNumber, authToken)
		if err != nil {
			// If there's an error, log the error and continue.
			log.Println(err.Error())
			continue
		}
		if err == nil && len(prs) == 0 {
			// no error but no data, that means no more pagination needed.
			break
		}
		// Filter the PRs based on teammates.
		prs = FilterList(prs, teammates)
		combinedPullRequests = append(combinedPullRequests, prs...)
	}

	return combinedPullRequests
}

// Main function.
func GetListOfPRs() (string, error) {
	var wg sync.WaitGroup
	repos := GetRepos()
	authToken := GetAuthToken()
	wg.Add(len(repos))

	maxConcurrent := 3 // max number of concurrent processes to spawn.
	sem := make(chan struct{}, maxConcurrent)

	var results []PullRequest

	// For each repo, we spawn a separate process and then combine the result from each process.
	for _, repo := range repos {
		sem <- struct{}{}
		go func(repo string) {
			defer func() { <-sem }()
			result := getPullRequests(repo, authToken)
			results = append(results, result...)
			wg.Done()
		}(repo)
	}
	wg.Wait()

	// Sort combined pull requests based on the PR age (oldest one comes first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Created_At < results[j].Created_At
	})

	// Truncate results if necessary.
	if len(results) > MaxResults {
		results = results[:MaxResults]
	}

	// Format response body for Slack
	var responseBody string
	for _, pr := range results {
		days := ConvertTimeToDay(pr.Created_At)
		responseBody += fmt.Sprintf("* %s <%s> %s Created %d day(s) ago\n", pr.Title, pr.URL, pr.User.Login, days)
	}

	return string(responseBody), nil
}
