/*
Copyright © 2025 NAME HERE spicyn.v2001@gmail.com
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			reader := bufio.NewReader(os.Stdin)

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
	users, _, err := git.Users.ListUsers(&gitlab.ListUsersOptions{})
	fmt.Printf("Found %d users", len(users))
}
