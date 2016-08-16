package scanner

import (
	"github.com/tjgq/sane"
)

type SaneOptions map[string][]sane.Option

func (s *Scanner) Options() SaneOptions {
	d, err := sane.Open(s.device.Name)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	options := make(map[string][]sane.Option)

	for _, o := range d.Options() {
		if o.IsSettable {
			options[o.Group] = append(options[o.Group], o)
		}
	}

	return options
}
