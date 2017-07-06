package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net"

	"github.com/elevran/relay/log"
)

const (
	defaultAddress = ":8080"
)

func main() {

	logger := log.NewStdlibLogger("relayd")

	address := flag.String("listen", defaultAddress, "Address:port on which to listen.")
	flag.Parse()

	if *address != defaultAddress {
		_, _, err := net.SplitHostPort(*address)
		if err != nil {
			logger.Fatal("Can't listen on", *address, err)
		}
	}

	relayServer := newRelay(*address, logger)

	go func() {
		logger.Println("Starting relay on", *address)
		err := relayServer.start()
		if err != nil {
			logger.Fatal("Failed to start relay", err)
		}
	}()

	// Wait for a shutdown signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// Shutdown the service and give it 30 seconds for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := relayServer.shutdown(ctx)
	if err != nil {
		logger.Fatal("Failed to shutdown relay frontend: ", err)
	}
}
