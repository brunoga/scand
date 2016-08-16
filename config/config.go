package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	m map[string]string
}

func New(configFilePath string) (*Config, error) {
	f, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[string]string)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		strippedText := strings.TrimSpace(scanner.Text())
		if len(strippedText) == 0 || strippedText[0] == '#' {
			// Ignore empty and comment lines.
			continue
		}

		keyValue := strings.SplitN(scanner.Text(), ":", 2)
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid config line : %q",
				scanner.Text())
		}

		m[strings.TrimSpace(keyValue[0])] =
			strings.TrimSpace(keyValue[1])
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Config{
		m,
	}, nil
}

func (c *Config) Validate(keys ...string) error {
	for _, key := range keys {
		_, ok := c.m[key]
		if !ok {
			return fmt.Errorf("required configuration option %q "+
				"not found", key)
		}
	}

	return nil
}

func (c *Config) Get(key string) string {
	return c.m[key]
}
