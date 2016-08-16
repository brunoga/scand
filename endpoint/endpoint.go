package endpoint

import (
	"crypto/md5"
	"fmt"
	"log"
	"sync"
	"time"

	"bga/scand/config"
	"bga/scand/scanner"
)

type endpoint struct {
	name string
	uid  string

	m        sync.Mutex
	instance string

	c *config.Config

	s scanner.Scanner
}

func New(name string, c *config.Config, s scanner.Scanner) *endpoint {
	uid := generateUid(name)

	return &endpoint{
		name,
		uid,
		sync.Mutex{},
		"",
		c,
		s,
	}
}

func generateUid(name string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(name)))
}

func (e *endpoint) Run(stopChannel <-chan struct{}, wg *sync.WaitGroup) error {
	log.Printf("Running endpoint %q (%s) for %q.", e.name, e.uid,
		e.s.Model())

	err := e.register()
	if err != nil {
		return err
	}

	wg.Add(1)
	go e.startWork(stopChannel, wg)

	return nil
}

func (e *endpoint) startWork(stopChannel <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	// Start timers for SMNP checks and registration refresh.
	smnpTimerChannel := time.After(500 * time.Millisecond)
	registrationTimerChannel := time.After(10 * time.Minute)

L:
	for {
		select {
		case <-stopChannel:
			log.Printf("Stopping endpoint %q (%s) for  %q.",
				e.name, e.uid, e.s.Model())
			e.unregister()
			break L
		case <-smnpTimerChannel:
			e.handleSnmpStatus(e.checkSnmpStatus())

			// reset timer.
			smnpTimerChannel = time.After(500 * time.Millisecond)
		case <-registrationTimerChannel:
			// Refresh registration.
			e.register()

			// reset timer.
			registrationTimerChannel = time.After(10 * time.Minute)
		}
	}
}
