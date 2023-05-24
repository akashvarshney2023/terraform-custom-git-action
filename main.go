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
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current working directory: %v\n", err)
		os.Exit(1)
	}

	modules, err := getModuleList(workingDir)
	if err != nil {
		fmt.Printf("Error getting module list: %v\n", err)
		os.Exit(1)
	}

	validateTags(modules)
}

func validateTags(modules map[string]string) {
	for moduleName, moduleSource := range modules {
		repoURL, tag, err := parseRepoURLAndTag(moduleSource)
		if err != nil {
			fmt.Printf("Error parsing repo URL and tag: %v\n", err)
			continue
		}

		latestTag, hasLatest, err := hasLatestTag(repoURL, tag)
		if err != nil {
			fmt.Printf("Error checking latest tag: %v\n", err)
			continue
		}

		if !hasLatest {
			fmt.Printf("\033[33mWarning: The module \033[32m%s\033[33m is not the latest version. Please consider using the latest tag, which is \033[32m%s\033[33m\n", moduleName, latestTag)

		}
	}
}

func getModuleList(workingDir string) (map[string]string, error) {
	config, err := tfconfig.LoadModule(workingDir)
	if err != nil {
		return nil, fmt.Errorf("error loading Terraform configuration: %s", err)
	}

	moduleList := make(map[string]string)
	for moduleName, module := range config.ModuleCalls {
		moduleList[moduleName] = module.Source
	}
	return moduleList, nil
}

func parseRepoURLAndTag(input string) (string, string, error) {
	input = strings.TrimPrefix(input, "git::")

	u, err := url.Parse(input)
	if err != nil {
		return "", "", err
	}

	repoURL := u.Scheme + "://" + u.Host + u.Path
	tag := u.Query().Get("ref")

	return repoURL, tag, nil
}

func hasLatestTag(repoURL, tag string) (string, bool, error) {
	parts := strings.Split(repoURL, "/")
	owner := parts[len(parts)-2]
	repoWithGit := parts[len(parts)-1]
	repo := strings.TrimSuffix(repoWithGit, ".git")
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return tag, false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tag, false, err
	}

	var tags []Tag
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return tag, false, err
	}

	latestTag := ""
	for _, t := range tags {
		if latestTag == "" || t.Name > latestTag {
			latestTag = t.Name
		}
	}

	if latestTag == tag {
		return latestTag, true, nil
	}

	return latestTag, false, nil
}

type Tag struct {
	Name string `json:"name"`
}
