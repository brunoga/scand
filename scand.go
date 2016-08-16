package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/brunoga/scand/config"
	"github.com/brunoga/scand/endpoint"
	"github.com/brunoga/scand/scanner"
)

var (
	homeDir        = os.Getenv("HOME")
	configFilePath = flag.String("config_file_path", homeDir+"/.scand",
		"path to configuration file")
)

func setupSignal() <-chan struct{} {
	// Create stop channel, used to signal goroutines to stop.
	stopChannel := make(chan struct{})

	// Setup signal handling channel.
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Start signal goroutine.
	go func() {
		sig := <-sigChannel

		log.Printf("Got signal %s. Stopping workers.", sig)

		// Signal all goroutines to stop by closing stopChannel.
		close(stopChannel)
	}()

	return stopChannel
}

func main() {
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c, err := config.New(*configFilePath)
	if err != nil {
		log.Fatalf("Error loading configuration file %q : %s\n",
			*configFilePath, err)
	}

	err = c.Validate(
		"EndpointEmails",
		"MessageSubject",
		"MessageBody",
		"MessageFromName",
		"MessageFromAddress",
		"SmtpAuthUser",
		"SmtpAuthPassword",
		"SmtpServerPort")
	if err != nil {
		log.Fatal("Error validating configuration file %q : %s\n",
			*configFilePath, err)
	}

	scanners, err := scanner.Detect()
	if err != nil {
		log.Fatal(err)
	}

	if len(scanners) == 0 {
		log.Fatal("No Samsung network scanners found.")
	}

	// Setup signal trapping.
	stopChannel := setupSignal()

	var wg sync.WaitGroup

	for _, s := range scanners {
		for _, endpointEmail := range strings.Fields(
			c.Get("EndpointEmails")) {
			// Create endpoint for scanner.
			e := endpoint.New(endpointEmail, c, s)

			// Run endpoint.
			err := e.Run(stopChannel, &wg)
			if err != nil {
				panic(err)
			}
		}
	}

	// Wait for all goroutines to exit.
	wg.Wait()
}
