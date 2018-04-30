package cmd

import (
	"github.com/jlucasnsilva/atog/stalker"
	"github.com/jlucasnsilva/atog/ui"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Watch the log files",
	Long:  `Pass to this command a list of files to be watched.`,
	Run: func(cmd *cobra.Command, args []string) {
		values := stalker.Watch(args)
		ui.Execute(args, values)
	},
}
