package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of atog",
	Long:  `Print the version number of atog`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Atog v0.1")
	},
}
