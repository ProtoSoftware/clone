package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

func main() {

	// Fetch env variables
	sourceRepo := os.Getenv("SOURCE_REPO")
	targetRepo := os.Getenv("TARGET_REPO")
	token := os.Getenv("GITHUB_TOKEN")

	// Check if the env variables are set
	if sourceRepo == "" || targetRepo == "" {
		panic("SOURCE_REPO and TARGET_REPO were not found in environment")
	}

	if token == "" {
		panic("GITHUB_TOKEN was not found in environment")
	}

	// Split the repo names
	sourceRepoSplit := strings.Split(sourceRepo, "/")
	targetRepoSplit := strings.Split(targetRepo, "/")

	// Check they are formatted correctly
	if len(sourceRepoSplit) != 2 || len(targetRepoSplit) != 2 {
		panic("SOURCE_REPO and TARGET_REPO must be in the format: user/repo")
	}

	// Create a new client
	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	client := github.NewClient(tc)

	// Fetch the source repo releases
	sourceReleases, _, err := client.Repositories.ListReleases(context.Background(), sourceRepoSplit[0], sourceRepoSplit[1], nil)
	if err != nil {
		panic(err)
	}

	// Fetch the target repo releases
	targetReleases, _, err := client.Repositories.ListReleases(context.Background(), targetRepoSplit[0], targetRepoSplit[1], nil)
	if err != nil {
		panic(err)
	}

	// Loop through and fine missing releases
	var missingReleases []*github.RepositoryRelease
	for _, targetRelease := range targetReleases {
		found := false
		for _, sourceRelease := range sourceReleases {
			if targetRelease.GetTagName() == sourceRelease.GetTagName() {
				found = true
				break
			}
		}
		if !found {
			missingReleases = append(missingReleases, targetRelease)
		}
	}

	var missing []string
	var skipped []string
	var created []string
	for _, missingRelease := range missingReleases {

		// SKip if it's a prerelease
		if missingRelease.GetDraft() {
			continue
		}

		missing = append(missing, missingRelease.GetTagName())

		// Skip releases over 3 month old (Prevents backlog of releases to process)
		if os.Getenv("SKIP_OLD") == "true" {
			if missingRelease.GetPublishedAt().AddDate(0, 1, 0).Before(time.Now()) {
				skipped = append(skipped, missingRelease.GetTagName())
				continue
			}
		}

		TagName := missingRelease.GetTagName()
		ReleaseName := missingRelease.GetName()
		ReleaseBody := fmt.Sprintf(`
<h3> Notice </h3>

Release automatically created from https://github.com/%s/%s using [clone](https://github.com/ProtoSoftware/clone)
		
<h1> Original Release Notes </h1>

%s
		`, targetRepoSplit[0], targetRepoSplit[1], missingRelease.GetBody())
		ReleaseAssets := missingRelease.Assets

		release, _, err := client.Repositories.CreateRelease(context.Background(), sourceRepoSplit[0], sourceRepoSplit[1], &github.RepositoryRelease{
			TagName: &TagName,
			Name:    &ReleaseName,
			Body:    &ReleaseBody,
			Assets:  ReleaseAssets,
		})

		if err != nil {
			panic(err)
		}

		created = append(created, release.GetTagName())
	}

	fmt.Println("Missing: " + strings.Join(missing, ", "))
	fmt.Println("Skipped: " + strings.Join(skipped, ", "))
	fmt.Println("Created: " + strings.Join(created, ", "))

	fmt.Println("Done! Added " + fmt.Sprintf("%d", len(created)) + " new releases 🎉")
}
