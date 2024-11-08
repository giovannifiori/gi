/*
Copyright © 2024 giovannifiori <gf@gfiori.dev>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gi",
	Short: "Quickly generate .gitignore files with one command",
	Long:  `gi is a CLI tool that saves you time by quickly generating .gitignore files for your projects with sane defaults from the Toptal API.`,
	Args:  cobra.ExactArgs(1),
	Run:   generateGitIgnore,
}

func generateGitIgnore(cmd *cobra.Command, args []string) {
	subject := args[0]

	resp, err := http.Get(fmt.Sprintf("https://www.toptal.com/developers/gitignore/api/%s", subject))
	if err != nil {
		fmt.Println("Error fetching .gitignore file from Toptal API", err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Printf("No .gitignore file found for %s\n", subject)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	fileContents := string(body)
	fmt.Println(fileContents)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
