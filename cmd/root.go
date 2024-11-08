/*
Copyright © 2024 giovannifiori <gf@gfiori.dev>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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
	if _, err := os.Stat(".gitignore"); err == nil {
		fmt.Printf("A .gitignore file already exists in this directory. Overwrite it? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.Trim(strings.Replace(response, "\n", "", -1), " "))

		for response != "y" && response != "n" && response != "Y" && response != "N" && response != "" {
			fmt.Printf("A .gitignore file already exists in this directory. Overwrite it? [y/N]: ")
			response, _ = reader.ReadString('\n')
			response = strings.ToLower(strings.Trim(strings.Replace(response, "\n", "", -1), " "))
		}

		if response != "y" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}
	}

	subject := args[0]

	resp, err := http.Get(fmt.Sprintf("https://www.toptal.com/developers/gitignore/api/%s", subject))
	if err != nil {
		fmt.Println("Error fetching .gitignore file from API", err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Printf("No .gitignore file found for %s\n", subject)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading .gitignore file from API", err)
		os.Exit(3)
	}

	err = os.WriteFile(".gitignore", body, 0644)
	if err != nil {
		fmt.Println("Error writing .gitignore file", err)
		os.Exit(4)
	}

	fmt.Printf("Generated .gitignore file for %s", subject)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
