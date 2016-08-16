package endpoint

import (
	"fmt"
	"log"
	"strings"

	"github.com/alouca/gosnmp"
)

func (e *endpoint) checkSnmpStatus() byte {
	s, err := gosnmp.NewGoSNMP(e.s.IP().String(), "public",
		gosnmp.Version2c, 5)
	if err != nil {
		panic(err)
	}

	e.m.Lock()

	if len(e.instance) == 0 {
		e.m.Unlock()
		return 0
	}

	resp, err := s.Get("1.3.6.1.4.1.236.11.5.11.81.11.7.2.1.2." +
		e.instance)
	if err != nil {
		if strings.HasSuffix(err.Error(), "i/o timeout\n") {
			// Timeout.
			e.m.Unlock()
			return 0
		}

		panic(fmt.Sprintf("%q", err))
	}

	e.m.Unlock()

	for _, v := range resp.Variables {
		// Return the first byte of the first value we get.
		return byte(v.Value.(string)[0])
	}

	return 0
}

func (e *endpoint) handleSnmpStatus(status byte) {
	switch status {
	case 0:
		// Do nothing.
	case 1:
		// Endpoint selected.
		// Send scan options.
		log.Printf("%s %q Endpoint selected. Sending scan options.\n",
			e.uid, e.name)
		_, err := e.sendScanOptions()
		if err != nil {
			panic(err)
		}
	case 2:
		// Scan options selected.
		// Check user selected options.
		// Execute scan.
		// Re-register endpoint.
		log.Printf("%s %q User scan options received.\n",
			e.uid, e.name)
		_, err := e.getUserScanOptions()
		if err != nil {
			panic(err)
		}

		// TODO(bga): Actually use the user scan options.

		log.Printf("%s %q Starting scan.\n", e.uid, e.name)
		data, err := e.s.Scan()
		if err != nil {
			panic(err)
		}

		log.Printf("%s %q Scan done.\n", e.uid, e.name)

		err = e.sendEmail(data)
		if err != nil {
			panic(err)
		}

	default:
		panic(fmt.Sprintf("Unhandled status %d.", status))
	}
}
