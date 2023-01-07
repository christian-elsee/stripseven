/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"

)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the stripseven reverse proxy",
	Long: `Run the stripseven reverse proxy`,
	Run: func(cmd *cobra.Command, args []string) {
   log.SetFormatter(&log.JSONFormatter{})
    log.WithFields(
        log.Fields{
            "foo": "foo",
            "bar": "bar",
        },
    ).Info("Something happened")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	runCmd.Flags().IntP("port", "p", 8080, "specify local bind port")
}
