# github-open-pull-requests
This repo hosts the code for generating a list of open pull requests filtered by a list of team members.
The code is used with an AWS Lambda, and the github token is stored in AWS Secrets Manager. For local testing, 
GITHUB_ACCESS_TOKEN environment variable can be set to bypass retrieving the secret from AWS Secrets Manager.

## Build commands for the Lambda function (as per https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html):
```
$ GOOS=linux GOARCH=amd64 go build main
$ zip main.zip main
```

## To build locally:

Add a main() function like the following.

```
func main() {
	os.Setenv("GITHUB_ACCESS_TOKEN", "<access_token>") // for testing purposes.
	response, err := GetListOfPRs()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		return
	}
	fmt.Println(response)
}
```

## Environment variables to be set:

```
AWS_REGION_NAME = a region name
[e.g. "us-east-1"]

AWS_SECRET_NAME = <Secret name stored in Secrets Manager>
[The secret name is the entry in AWS Secrets Manager which stores the github authorization token that is used during API calls]

NUM_PAGES = <Maximum pages to pull from a repo>
[Github API retrieves 100 pull requests at a time, referred as page. So if NUM_PAGES=7, then in total 700 pull requests will be retrieved.]

REPOS = <Repo names separated by comma>
[The repository names to be chcked through the API. For example, "repo1,repo2"]

TEAMMATES = <Github usernames separated by comma>
[The list of people whose pull requests are to be retrieved. For example, "johndoe,janerobinson"]

REPO_OWNER = <owner of the repositories>
```
