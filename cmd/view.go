package cmd

import "github.com/spf13/cobra"

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Watch the log files",
	Long:  `Pass to this command a list of files to be watched.`,
}

func init() {
	previous := uint(0)
	size := uint(0)

	viewCmd.Flags().UintVarP(&previous, "previous", "p", 10, "Number of lines read from the end of the files when it is openned.")
	viewCmd.Flags().UintVarP(&size, "size", "s", 50, "Number of messages it will keep in the circular buffer.")
}
