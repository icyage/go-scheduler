package cmd

import (
	"github.com/obsurvive/voyager/log"
	"github.com/obsurvive/voyager/scheduler"
	"github.com/spf13/cobra"
)

// schedulerCmd represents the version command
var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Starts the scheduler service",
	Long:  `Starts scheduler server`,
	Run: func(cmd *cobra.Command, args []string) {
		scheduler.Run()
	},
}

// schedulerCmd represents the version command
var schedulerInitDBCmd = &cobra.Command{
	Use:   "initdb",
	Short: "Creates needed tables",
	Long:  `Initializes the database for scheduling`,
	Run: func(cmd *cobra.Command, args []string) {
		err := scheduler.CreateTable()
		log.Error(err)
	},
}

func init() {
	schedulerCmd.AddCommand(schedulerInitDBCmd)
	rootCmd.AddCommand(schedulerCmd)
}
