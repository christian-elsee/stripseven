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

	ln, err := net.ListenTCP("tcp", lst)
	if err != nil {
		logf.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to bind to listen address")
		return err
	}
	logf.Info("Bind to listen address")

	conntgt, err := net.DialTCP("tcp", nil, tgt)
	if err != nil {
		logf.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to dial target address")
		return err
	}

	connlst, err := ln.Accept()
	if err != nil {
		logf.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to accept connection on listen adddress")
		return err
	}

	logf.Debug("Start two-way relay")
	wg.Add(2)

	go func() {
		logf.Info("Relay client to target service")
		defer wg.Done()

		relay(ctx, connlst, conntgt)
	}()
	go func() {
		logf.Info("Relay target service to client")
		defer wg.Done()

		relay(ctx, conntgt, connlst)
	}()

	logf.Info("Waiting on relay to close")
	wg.Wait()
	logf.Info("Finished waiting for relay to close")

	connlst.Close()
	logf.Debug("The listen connection is closed")

	conntgt.Close()
	logf.Debug("The target connection is closed")

	return nil
}

func relay(ctx context.Context, src, dst net.Conn) {
  logf := log.WithFields(
    log.Fields{
    	"trace": Trace("relay", "pkg/proxy"),
    	"src": src,
    	"dst": dst,
    },
  )
  logf.Debug("Enter")
  defer logf.Debug("Exit")

  scanner := bufio.NewScanner(bufio.NewReader(src))
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


