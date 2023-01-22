/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"net"
	"os"
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"

	"github.com/christianlc-highlights/stripseven/pkg"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the stripseven reverse proxy",
	Long: `Run the stripseven reverse proxy`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGHUP,
		)
    logf := log.WithFields(
      log.Fields{
      	"trace": pkg.Trace("runCmd.Run", "cmd/run"),
      	"listen": pkg.Must(cmd.Flags().GetString("listen")),
      	"target": pkg.Must(cmd.Flags().GetString("target")),
      },
    )

    logf.Debug("Enter")
    defer logf.Debug("Exit")
    defer stop()

    lst, err := net.ResolveTCPAddr("tcp", pkg.Must(cmd.Flags().GetString("listen")))
    if err != nil {
    	logf.
    		WithFields(log.Fields{ "err": err }).
    		Error("Failed to resolve listen address")
    	panic(err)
    }

    tgt, err := net.ResolveTCPAddr("tcp", pkg.Must(cmd.Flags().GetString("target")))
    if err != nil {
    	logf.
    		WithFields(log.Fields{ "err": err }).
    		Error("Failed to resolve target address")
    	panic(err)
    }

    logf.Info("Start proxy")
		err = pkg.Proxy(ctx, lst, tgt)
		if err != nil {
			logf.
		  	WithFields(log.Fields{ "err": err }).
		  	Error("Proxy has stopped")
		}
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

	runCmd.Flags().StringP("listen", "l", "localhost:8080", "specificy listen address")
	runCmd.Flags().StringP("target", "t", "", "specify target address")
	runCmd.MarkFlagRequired("target")
}
