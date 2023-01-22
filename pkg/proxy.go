package pkg

import (
	"net"
	"bufio"
	"sync"
	"fmt"
	"context"

	log "github.com/sirupsen/logrus"

)
// func /////////////////////////////////////////

// Establishes a persistent, two-way proxy between
// listen and target addresses
func Proxy(ctx context.Context, lst, tgt *net.TCPAddr) error {
	var wg sync.WaitGroup

  logf := log.WithFields(
    log.Fields{
    	"trace": Trace("Proxy", "pkg/proxy"),
    	"listen": lst,
    	"target": tgt,
    },
  )

  logf.Debug("Enter")
  defer logf.Debug("Exit")

	conntgt, err := net.DialTCP("tcp", nil, tgt)
	if err != nil {
		logf.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to dial target address")
		return err
	}
	logf.Info("Dialed target address")

	ln, err := net.ListenTCP("tcp", lst)
	if err != nil {
		logf.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to bind to listen address")
		return err
	}
	logf.Info("Bind to listen address")

	go func() {
		logf.Debug(
			"Wait on closed context to close relay connections",
		)
		<-ctx.Done()
		logf.Info("Close relay connections")
		defer logf.Info("Closed relay connections")

		ln.Close()
		conntgt.Close()

	}()

	connlst, err := ln.Accept()
	if err != nil {
		logf.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to accept connection on listen adddress")
		return err
	}

	logf.Info("Start relay")
	wg.Add(2)

	go func() {
		logf.Info("Relay client to target service")
		defer logf.Info("Finished client to target relay")
		defer wg.Done()

		relay(ctx, connlst, conntgt)
	}()
	go func() {
		logf.Info("Relay target service to client")
		defer logf.Info("Finished target to client relay")
		defer wg.Done()

		relay(ctx, conntgt, connlst)
	}()

	logf.Info("Waiting on relay to close")
	wg.Wait()
	logf.Info("Finished waiting for relay to close")

	return nil
}

func relay(ctx context.Context, src, dst net.Conn) {
  logf := log.WithFields(
    log.Fields{
    	"trace": Trace("relay", "pkg/proxy"),
    	"src": src.LocalAddr().String(),
    	"dst": dst.LocalAddr().String(),
    },
  )
  logf.Debug("Enter")
  defer logf.Debug("Exit")

  scanner := bufio.NewScanner(src)
  scanner.Split(bufio.ScanBytes)

  done := false
  for !done {

  	select {
  	case <-ctx.Done():
	  	logf.Info("Context has signaled to close relay")
	  	done = true
	  default:
	  	logf.Info("Read from src connnection")
	  	scanner.Scan()
	  	fmt.Println(scanner.Text())
  	}
  }

}


