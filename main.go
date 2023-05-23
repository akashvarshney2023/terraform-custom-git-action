package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

func main() {
	// Readcheck the Terraform configuration file
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current working directory: %v\n", err)
		os.Exit(1)
	}
	modules := getModuleList(workingDir)
	// Validate module tags
	validateTags(modules)
}

// Verify if the module is the latest version based on its tags
func validateTags(modules map[string]string) {

	for moduleName, moduleSource := range modules {
		repoURL, tag, err := ParseRepoURLAndTag(moduleSource)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		latestTag, hasLatest, err := HasLatestTag(repoURL, tag)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if !hasLatest {
			fmt.Printf("\033[33mWarning: The \033[36m %s \033[36m \033[33m is not the latest version. Please consider using the latest tag, which is\033[33m \033[36m %s\033[36m\n", moduleName, latestTag)
		}
	}

}

// Get the module list and name on the given specific path
func getModuleList(workingDir string) map[string]string {
	config, err := tfconfig.LoadModule(workingDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading Terraform configuration: %s/n", err)
		os.Exit(1)
	}

	moduleList := make(map[string]string)
	for moduleName, module := range config.ModuleCalls {
		moduleList[moduleName] = module.Source
	}
	return moduleList
}

func ParseRepoURLAndTag(input string) (string, string, error) {
	// Remove the "git::" prefix if present
	input = strings.TrimPrefix(input, "git::")

	// Parse the input as a URL
	u, err := url.Parse(input)
	if err != nil {
		return "", "", err
	}

	// Extract the repository URL
	repoURL := u.Scheme + "://" + u.Host + u.Path

	// Extract the tag from the query parameters
	tag := u.Query().Get("ref")

	return repoURL, tag, nil
}

func HasLatestTag(repoURL, tag string) (string, bool, error) {
	// Extract the owner and repository name from the URL
	parts := strings.Split(repoURL, "/")
	owner := parts[len(parts)-2]
	repoWithGit := parts[len(parts)-1]
	// Remove the ".git" suffix from the repository name
	repo := strings.TrimSuffix(repoWithGit, ".git")
	// Construct the GitHub API URL for tags
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo)

	// Send an HTTP GET request to the API
	resp, err := http.Get(apiURL)
	if err != nil {
		return tag, false, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tag, false, err
	}

	// Parse the response body as an array of tags
	var tags []Tag
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return tag, false, err
	}

	// Find the latest tag
	latestTag := ""
	for _, t := range tags {
		if latestTag == "" || t.Name > latestTag {
			latestTag = t.Name
		}
	}

	// Compare the latest tag with the provided tag
	if latestTag == tag {
		return latestTag, true, nil
	}

	return latestTag, false, nil
}

type Tag struct {
	Name string `json:"name"`
}
