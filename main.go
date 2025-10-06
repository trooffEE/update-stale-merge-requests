/*
Copyright © 2025 NAME HERE spicyn.v2001@gmail.com
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	"golang.org/x/term"
)

func main() {
	viper.SetConfigName("update-stale-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/")
	reader := bufio.NewReader(os.Stdin)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			fmt.Print("Enter your gitlab host (for example \"gitlab.com\"): ")
			host, _ := reader.ReadString('\n')
			host = strings.TrimSpace(host)

			fmt.Printf(
				"Paster your API token here (should be created here with api scope - https://%s/-/user_settings/personal_access_tokens): ",
				host,
			)
			secret, err := term.ReadPassword(syscall.Stdin)

			if err != nil {
				fmt.Printf("Error reading api token: %v\n", err)
			}

			viper.Set("gitlab.host", host)
			viper.Set("gitlab.token", string(secret))

			err = viper.SafeWriteConfig()
		} else {
			panic(err)
		}
	}

	git, err := gitlab.NewClient(
		viper.Get("gitlab.token").(string),
		gitlab.WithBaseURL(fmt.Sprintf("https://%s/api/v4", viper.Get("gitlab.host").(string))),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	fmt.Println("Run time GitLab client created")

	currentUser, _, err := git.Users.CurrentUser()
	if err != nil {
		log.Fatalf("Failed to get information about owner of token: %v", err)
	}
	fmt.Printf("Actions will be performed with rights of the owner of token - \"%s\" (%s)\n", currentUser.Username, currentUser.Email)

	fmt.Printf("List of available projects: \n")
	projects, _, err := git.Projects.ListProjects(&gitlab.ListProjectsOptions{
		OrderBy: gitlab.Ptr("updated_at"),
	})
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}
	for index, p := range projects {
		fmt.Printf("%d. %s\n", index+1, p.Name)
	}

	fmt.Println("Choose one project (provide order number): ")
	pString, _ := reader.ReadString('\n')
	pString = strings.TrimSpace(pString)
	projectOrder, err := strconv.Atoi(strings.TrimSpace(pString))
	if err != nil {
		log.Fatalf("Error converting project id to int: %v\n", err)
	}
	projectIndex := projectOrder - 1
	if projects[projectIndex] == nil {
		log.Fatalf("Project does not exist: %v", pString)
	}

	fmt.Printf("Searching MRs for project ID %d 🔎...", projects[projectIndex].ID)
	mergeRequests, _, err := git.MergeRequests.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
		State:    gitlab.Ptr("opened"),
		AuthorID: gitlab.Ptr(currentUser.ID),
		Scope:    gitlab.Ptr("assigned_to_me"),
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(mergeRequests) == 0 {
		log.Fatal("No open merge requests found 🚫")
	}

	for _, mr := range mergeRequests {
		if mr.ProjectID == projects[projectIndex].ID {
			fmt.Printf("\n%s (%s)", mr.Title, strings.Join(mr.Labels, ", "))
		}
	}
}
