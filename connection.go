package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
)

// Connection data
type Connection struct {
	NetworkType string
	URL         *url.URL
}

var defaultProtPorts = map[string]string{
	"http":  "80",
	"https": "443",
	"ssh":   "22",
}

// BuildConn build a connection structure.
// This connection data can later be used as a common structure
// by the functions that will check if the target is available.
func BuildConn(cfg *Config) (*Connection, error) { // nolint gocyclo
	address := cfg.Address
	if address == "" && cfg.Host != "" {
		address = net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

		if cfg.Protocol != "" {
			address = fmt.Sprintf("%s://%s", cfg.Protocol, address)
		}
	}

	if address == "" {
		return nil, errors.New("Connection address is empty")
	}

	u, err := url.Parse(address)

	// the url parsing may fail if it is missing the scheme value
	// so it try again adding a tcp scheme as falback
	var err2 error
	if err != nil || (u.Scheme != "" && u.Host == "") {
		u, err = url.Parse(fmt.Sprintf("tcp://%s", address))
	}

	// return error from the original address
	if err2 != nil {
		return nil, fmt.Errorf("Error parsing connection address: %v", err)
	}

	if u == nil || u.Hostname() == "" {
		return nil, fmt.Errorf("Couldn't parse address: %s", address)
	}

	p := u.Port()

	// resolve default port based on the provided scheme
	if p == "0" {
		if dp, ok := defaultProtPorts[u.Scheme]; ok {
			var re = regexp.MustCompile(`:0$`)
			u.Host = re.ReplaceAllString(u.Host, ":"+dp)
		}
	}

	// resolve default scheme based on the provided port
	if u.Scheme == "" {
		if p == "80" {
			u.Scheme = "http"
		} else if p == "443" {
			u.Scheme = "https"
		}
	}

	return &Connection{
		NetworkType: "tcp",
		URL:         u,
	}, nil
}
