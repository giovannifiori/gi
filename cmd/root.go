/*
Copyright © 2024 giovannifiori <gf@gfiori.dev>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/erikgeiser/promptkit/selection"
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
	cobra.CheckErr(err)
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Printf("No .gitignore file found for %s\n", subject)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Failed to fetch .gitignore file for %s\n", subject)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	cobra.CheckErr(err)

	if _, err := os.Stat(".gitignore"); err == nil {
		sp := selection.New("A .gitignore file already exists in this directory. Overwrite or append to it?", []string{"Append", "Overwrite"})
		sp.Filter = nil
		choice, err := sp.RunPrompt()
		cobra.CheckErr(err)

		if choice == "Append" {
			f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, 0644)
			cobra.CheckErr(err)
			defer f.Close()

			_, err = f.Write(body)
			cobra.CheckErr(err)

			fmt.Printf("Appended to .gitignore file with the contents for %s\n", subject)
		} else {
			err = os.WriteFile(".gitignore", body, 0644)
			cobra.CheckErr(err)

			fmt.Printf("Generated .gitignore file for %s\n", subject)
		}

	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
