package cmd

import (
	"github.com/obsurvive/voyager/http"
	"github.com/spf13/cobra"
)

// serveCmd represents the version command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts http server for the api",
	Long:  `Voyager start server`,
	Run: func(cmd *cobra.Command, args []string) {
		http.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
