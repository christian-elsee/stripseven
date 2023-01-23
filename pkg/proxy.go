package pkg

import (
	"io"
	"net"
	"fmt"
	"sync"
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

)
// func /////////////////////////////////////////

// Establishes a persistent, two-way proxy between
// listen and target addresses
func Proxy(ctx context.Context, lst, tgt *net.TCPAddr) error {
  ctx, cancel := context.WithCancel(ctx)
  logf := log.WithFields(
    log.Fields{
    	"trace": Trace("Proxy", "pkg/proxy"),
    	"listen": lst,
    	"target": tgt,
    },
  )
  logf.Debug("Enter")
  defer logf.Debug("Exit")
  defer cancel()

  server := http.Server{
  	Addr: lst.String(),
  }
  logf.Info("Create listen server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		logf.Info("Start relay request")
		defer logf.Debug("Finished relay request")

		conntgt, err := net.DialTCP("tcp", nil, tgt)
		if err != nil {
			logf.WithFields(log.Fields{
				"err": err,
			}).Error("Failed to dial target address")

			cancel()
			return
		}
		defer conntgt.Close()
		logf.Info("Dialed target address")

		// create two-way relay, first request body to target service
		wg.Add(2)
		go func() {
			logf.Debug("Relay request to target")
			defer wg.Done()

			_, err := io.Copy(conntgt, r.Body)
			logf.WithFields(log.Fields{
				"err": err,
			}).Debug("Stopped relay request to target")
		}()

		// and target service to response
		go func() {
			logf.Debug("Relay target to response")
			defer wg.Done()

			_, err := io.Copy(w, conntgt)
			fmt.Fprintln(w, "")

			logf.WithFields(log.Fields{
				"err": err,
			}).Debug("Stopped relay target to response")
		}()

		wg.Wait()
	})

	go func() {
		logf.Info("Start listen server")

		err := server.ListenAndServe()
		logf.WithFields(log.Fields{
			"err": err,
		}).Info("Closed listen server")
	}()

	logf.Debug("Wait for shutdown signal")
	<-ctx.Done()
	logf.Info("Received shutdown signal")

	err := server.Shutdown(ctx)
	if err != nil {
		logf.WithFields(log.Fields{
			"err": err,
		}).Info("Listen server shutdown")
		server.Close()
	}

	return err
}
