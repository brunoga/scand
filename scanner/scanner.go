package scanner

import (
	"bytes"
	"image/jpeg"
	"log"
	"net"
	"strings"

	"github.com/tjgq/sane"
)

func init() {
	if err := sane.Init(); err != nil {
		panic(err)
	}
}

type Scanner struct {
	device sane.Device
	ip     net.IP
}

func Detect() ([]Scanner, error) {
	log.Print("Detecting scanners using SANE.")

	devices, err := sane.Devices()
	if err != nil {
		return nil, err
	}

	scanners := make([]Scanner, 0, len(devices))

	for _, device := range devices {
		if device.Type == "Scanner" {
			s := Scanner{
				device,
				nil,
			}

			if s.isNetScanner() && s.isSamsungScanner() {
				// Found a network Samsung scanner.

				// Get IP.
				stringIP := s.device.Name[strings.Index(
					s.device.Name, ";")+1:]
				s.ip = net.ParseIP(stringIP)

				scanners = append(scanners, s)
			}
		}
	}

	log.Printf("Found %d Samsung network scanner(s).", len(scanners))

	return scanners, nil
}

func (s *Scanner) IP() net.IP {
	return s.ip
}

func (s *Scanner) Model() string {
	return s.device.Model
}

func (s *Scanner) Scan() ([]byte, error) {
	c, err := sane.Open(s.device.Name)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// TODO(bga): Set options.

	img, err := c.ReadImage()
	if err != nil {
		return nil, err
	}

	var f bytes.Buffer

	err = jpeg.Encode(&f, img, nil)
	if err != nil {
		return nil, err
	}

	return f.Bytes(), nil
}

func (s *Scanner) isNetScanner() bool {
	return strings.HasPrefix(s.device.Name, "smfp:net;")
}

func (s *Scanner) isSamsungScanner() bool {
	return s.device.Vendor == "Samsung"
}
