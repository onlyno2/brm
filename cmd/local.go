/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/onlyno2/brm/utils"
	"github.com/spf13/cobra"
)

// flag definitions
var projectDir string

func LocalBranches() []string {
	branchesCmd := exec.Command(
		"git",
		"for-each-ref",
		"--sort=committerdate",
		"refs/heads/",
		"--format=(%(committerdate:relative))\t%(refname:short)\t%(authorname)")
	branchesCmd.Dir = projectDir
	out, err := branchesCmd.Output()
	utils.CheckErr(err)

	branches := strings.Split(string(out), "\n")
	return branches[:len(branches)-1]
}

func CheckBoxes(label string, branches []string) []string {
	res := []string{}

	prompt := &survey.MultiSelect{
		Message: label,
		Options: branches,
	}

	survey.AskOne(prompt, &res)

	return res
}

func DeleteBranches(branches []string) {
	for _, branch := range branches {
		fmt.Println("Deleting ", branch)
		deleteCmd := exec.Command("git", "branch", "-D", branch)
		deleteCmd.Dir = projectDir
		_, err := deleteCmd.Output()
		utils.CheckErr(err)
	}
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Manage git local branches",
	Long:  `Manage git local branches`,
	Run: func(cmd *cobra.Command, args []string) {
		selectedOptions := CheckBoxes("Select branches to delete:", LocalBranches())
		deleteBranches := make([]string, len(selectedOptions))
		for i, option := range selectedOptions {
			tokens := strings.Split(option, "\t")

			deleteBranches[i] = tokens[1]
		}

		DeleteBranches(deleteBranches)
		fmt.Println("OK")
	},
}

func init() {
	rootCmd.AddCommand(localCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	localCmd.Flags().StringVar(&projectDir, "dir", ".", "Specify directory (default is current directory)")
}
