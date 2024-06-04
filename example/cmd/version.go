/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = ""
	GitRev    = ""
	BuildTime = ""
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print app's version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nVersion:	 v%s\nbuildTime:	 %s\ngitRev:\t	 %s\n", Version, BuildTime, GitRev)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
