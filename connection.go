package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
)

// Connection data
type Connection struct {
	NetworkType string
	URL         *url.URL
}

// BuildConn build a connection structure.
// This connection data can later be used as a common structure
// by the functions that will check if the target is available.
func BuildConn(cfg *Config) (*Connection, error) { // nolint gocyclo
	address := cfg.Address
	if cfg.Host != "" {
		address = net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
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

	if u.Hostname() == "" {
		return nil, fmt.Errorf("Couldn't parse address: %s", address)
	}

	if u.Scheme == "" {
		if p := u.Port(); p == "80" {
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
