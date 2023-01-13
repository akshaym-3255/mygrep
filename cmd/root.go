package cmd

import (
	"fmt"
	"os"

	"github.com/akshaym-3255/mygrep/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mygrep",
	Short: "mygrep is command line tool to search  pattern",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println("pattern to search is required")
			return
		}
		pattern := args[0]

		insensitive, _ := cmd.Flags().GetBool("insensitive")
		recursive, _ := cmd.Flags().GetBool("recursive")
		if recursive && len(args) == 1 {
			fmt.Println("directory name is required")
			return
		}
		outPutFile, _ := cmd.Flags().GetString("outputFile")

		location := ""
		readFromStdIn := false

		if len(args) == 1 {
			readFromStdIn = true
		} else {
			location = args[1]
		}

		grepCommand := internal.Grep{
			CaseInSensitive: insensitive,
			Recursive:       recursive,
			ReadFromStdIn:   readFromStdIn,
			Path:            location,
			Pattern:         pattern,
			OutputFile:      outPutFile,
		}

		matchedLines, err := grepCommand.MatchPattern()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		grepCommand.WriteOutput(matchedLines)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	insensitive := false
	recursive := false
	outPutFile := ""
	rootCmd.Flags().BoolVarP(&insensitive, "insensitive", "i", false, "insensitive match")
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "recursive search")
	rootCmd.Flags().StringVarP(&outPutFile, "outputFile", "o", "", "write output to file")
}
