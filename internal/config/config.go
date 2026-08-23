package config

import (
	"os"
	"strconv"
	"strings"
)

// Config contains runtime settings. Values can be loaded from a tiny YAML subset
// and are always overridable with environment variables.
type Config struct {
	Host            string
	Port            int
	Environment     string
	ShutdownSeconds int
	LogLevel        string
}

func Default() Config {
	return Config{Host: "0.0.0.0", Port: 8092, Environment: "development", ShutdownSeconds: 10, LogLevel: "info"}
}

func Load(path string) (Config, error) {
	c := Default()
	if path != "" {
		if b, err := os.ReadFile(path); err == nil {
			parseYAML(string(b), &c)
		} else if !os.IsNotExist(err) {
			return c, err
		}
	}
	if v := os.Getenv("ALERT_HOST"); v != "" {
		c.Host = v
	}
	if v := os.Getenv("ALERT_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.Port = n
		}
	}
	if v := os.Getenv("ALERT_ENV"); v != "" {
		c.Environment = v
	}
	if v := os.Getenv("ALERT_SHUTDOWN_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.ShutdownSeconds = n
		}
	}
	if v := os.Getenv("ALERT_LOG_LEVEL"); v != "" {
		c.LogLevel = v
	}
	return c, nil
}

func parseYAML(s string, c *Config) {
	for _, line := range strings.Split(s, "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		switch key {
		case "host":
			c.Host = value
		case "port":
			if n, e := strconv.Atoi(value); e == nil {
				c.Port = n
			}
		case "environment":
			c.Environment = value
		case "shutdown_seconds":
			if n, e := strconv.Atoi(value); e == nil {
				c.ShutdownSeconds = n
			}
		case "log_level":
			c.LogLevel = value
		}
	}
}

func ParseYAML(s string, c *Config) { parseYAML(s, c) }
