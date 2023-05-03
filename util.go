package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var DefaultMaxNumberOfPages = 10 // max pages for each repo to retrieve
var MaxResults = 100             // maximum entries in the final result

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func SplitAndTrim(input string, separator string) []string {
	slc := strings.Split(input, separator)
	for i := range slc {
		slc[i] = strings.TrimSpace(slc[i])
	}
	return slc
}

func GetTeammates() []string {
	teammatesVar := os.Getenv("TEAMMATES")
	if len(teammatesVar) != 0 {
		// split the env variable.
		return SplitAndTrim(teammatesVar, ",")
	} else {
		fmt.Println("Comma-separated list of team-members not found in env variable.")
		return []string{}
	}
}

func GetRepos() []string {
	reposVar := os.Getenv("REPOS")
	if len(reposVar) != 0 {
		// split the env variable.
		return SplitAndTrim(reposVar, ",")
	} else {
		fmt.Println("Comma-separated list of repositories not found in env variable.")
		return []string{}
	}
}

func GetAuthToken() string {
	// First see if the env variable is defined.
	tokenVar := os.Getenv("GITHUB_ACCESS_TOKEN")
	if len(tokenVar) != 0 {
		fmt.Println("Token found in env variable.")
		return tokenVar
	} else {
		// Try from aws secrets manager.
		fmt.Println("Token not found in env variable, trying the alternative.")
		return GetSecret()
	}
}

func GetRepoOwner() string {
	owner := os.Getenv("REPO_OWNER")
	if len(owner) != 0 {
		return owner
	} else {
		fmt.Println("Repo owner not found in env variable.")
		return ""
	}
}

func GetSlackWebhookUrl() string {
	url := os.Getenv("SLACK_WEBHOOK_URL")
	if len(url) != 0 {
		return url
	} else {
		fmt.Println("Slack webhook URL not found in env variable.")
		return ""
	}
}

func GetAwsRegion() string {
	region := os.Getenv("AWS_REGION_NAME")
	if len(region) != 0 {
		return region
	} else {
		fmt.Println("AWS region not found in env variable.")
		return ""
	}
}

func GetAwsSecretName() string {
	secretName := os.Getenv("AWS_SECRET_NAME")
	if len(secretName) != 0 {
		return secretName
	} else {
		fmt.Println("AWS secret name not found in env variable.")
		return ""
	}
}

func GetNumberOfPages() int {
	nVar := os.Getenv("NUM_PAGES")
	n, _ := strconv.Atoi(nVar)
	if n != 0 {
		return n
	}
	return DefaultMaxNumberOfPages
}

// Filter an array of PRs to only contain PRs authored by someone in the teammates array.
func FilterList(pullRequests []PullRequest, teammates []string) []PullRequest {
	var filteredPullRequests []PullRequest
	for _, v := range pullRequests {
		if Contains(teammates, v.User.Login) {
			filteredPullRequests = append(filteredPullRequests, v)
		}
	}
	return filteredPullRequests
}

func ConvertTimeToDay(timestamp string) int {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		// in case of error.
		return 0
	}
	days := int(time.Since(t).Hours() / 24)
	return days
}
