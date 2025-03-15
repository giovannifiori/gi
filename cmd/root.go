/*
Copyright © 2024 giovannifiori <gf@gfiori.dev>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gi",
	Short: "Quickly generate .gitignore files with one command",
	Long:  `gi is a CLI tool that saves you time by quickly generating .gitignore files for your projects with sane defaults from the Toptal API.`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	subjects, err := getSubjects(args)
	cobra.CheckErr(err)

	if len(subjects) == 0 {
		fmt.Println("No subject selected. Exiting...")
		os.Exit(0)
	}

	formattedSubjects := strings.Join(subjects, ",")

	fileContents, err := getFileContents(subjects)
	cobra.CheckErr(err)

	if _, err := os.Stat(".gitignore"); err == nil {
		var choice string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("A .gitignore file already exists in this directory. Overwrite or append to it?").
					Options(huh.NewOption("Append", "append"), huh.NewOption("Overwrite", "overwrite")).
					Value(&choice),
			),
		)

		err := form.Run()
		cobra.CheckErr(err)

		if choice == "append" {
			err = appendToGitIgnoreFile(fileContents)
			cobra.CheckErr(err)
			fmt.Printf("Appended to .gitignore file with the content for %s\n", formattedSubjects)
		} else {
			err = writeGitIgnoreFile(fileContents)
			cobra.CheckErr(err)
			fmt.Printf("Overwritten .gitignore file with the content for %s\n", formattedSubjects)
		}
	} else {
		err = writeGitIgnoreFile(fileContents)
		cobra.CheckErr(err)
		fmt.Printf("Generated .gitignore file for %s\n", formattedSubjects)
	}
}

func getFileContents(subjects []string) (body []byte, err error) {
	parsedSubjects := strings.Join(subjects, ",")
	resp, err := http.Get(fmt.Sprintf("https://www.toptal.com/developers/gitignore/api/%s", parsedSubjects))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("no .gitignore content found for %s", parsedSubjects)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch .gitignore file for %s", parsedSubjects)
	}

	return io.ReadAll(resp.Body)
}

func getSubjects(args []string) ([]string, error) {
	var subjects []string

	if len(args) > 0 {
		return args, nil
	}

	subjectForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Generate .gitignore for: ").
				Value(&subjects).
				OptionsFunc(func() []huh.Option[string] {
					resp, err := http.Get("https://www.toptal.com/developers/gitignore/api/list")
					cobra.CheckErr(err)
					defer resp.Body.Close()

					if resp.StatusCode != 200 {
						fmt.Printf("Failed to list options")
						os.Exit(1)
					}

					optsBodyBytes, optsBodyErr := io.ReadAll(resp.Body)
					cobra.CheckErr(optsBodyErr)

					optsStr := string(optsBodyBytes)

					opts := strings.Split(strings.Join(strings.Split(optsStr, "\n"), ","), ",")

					return huh.NewOptions(opts...)
				}, nil),
		),
	)

	err := subjectForm.Run()
	if err != nil {
		return nil, err
	}

	return subjects, nil
}

func writeGitIgnoreFile(data []byte) error {
	err := os.WriteFile(".gitignore", data, 0644)
	return err
}

func appendToGitIgnoreFile(data []byte) error {
	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
