package cmd

import (
	"fmt"
	"github.com/obsurvive/voyager/log"
	"github.com/obsurvive/voyager/probes"
	"github.com/spf13/cobra"
)

var (
	URL      string
	Insecure bool
)

// serveCmd represents the version command
var probeCmd = &cobra.Command{
	Use:   "probe",
	Short: "probe url",
	Long:  `Probes a given url`,
	Run: func(cmd *cobra.Command, args []string) {
		url := probes.Check{}
		url.URL = URL
		url.Insecure = Insecure

		response, err := probes.CheckHTTP(&url)
		if err != nil {
			log.Error(err)
		}
		fmt.Printf("%+v\n", response)
		fmt.Printf("%+v\n", response.Timeline)
		fmt.Printf("%+v\n", response.SSL)
	},
}

var probeLauncherCmd = &cobra.Command{
	Use:   "launcher",
	Short: "starts launch control for probes",
	Long:  `reads from queue and makes sure probes launches correctly`,
	Run: func(cmd *cobra.Command, args []string) {
		probes.Launch()
	},
}

var probeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "gets probe",
	Long:  `reads probe details`,
	Run: func(cmd *cobra.Command, args []string) {
		probe, err := probes.GetProbe(args[0])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(probe)
	},
}

func init() {
	rootCmd.AddCommand(probeCmd)
	probeCmd.Flags().StringVarP(&URL, "url", "u", "", "url to probe")
	probeCmd.Flags().BoolVarP(&Insecure, "insecure", "i", false, "should we validate ssl?")
	probeCmd.MarkFlagRequired("url")
	probeCmd.AddCommand(probeLauncherCmd)
	probeCmd.AddCommand(probeGetCmd)
}
