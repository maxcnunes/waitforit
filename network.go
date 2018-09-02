package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

// DialConfigs dial multiple connections at same time
func DialConfigs(confs []Config, print func(a ...interface{})) error {
	ch := make(chan error)
	for _, config := range confs {
		go func(conf Config) {
			conn, err := BuildConn(&conf)
			if err != nil {
				ch <- fmt.Errorf("Invalid connection %#v: %v", conf, err)
				return
			}

			ch <- DialConn(conn, &conf, print)
		}(config)
	}

	for i := 0; i < len(confs); i++ {
		if err := <-ch; err != nil {
			return err
		}
	}

	return nil
}

// DialConn check if the connection is available
func DialConn(conn *Connection, conf *Config, print func(a ...interface{})) error {
	print("Waiting " + strconv.Itoa(conf.Timeout) + " seconds")
	if err := pingHost(conn, conf, print); err != nil {
		return err
	}

	if conn.URL.Scheme == "http" || conn.URL.Scheme == "https" {
		return pingAddress(conn, conf, print)
	}

	return nil
}

// pingAddress check if the full address is responding properly
func pingAddress(conn *Connection, conf *Config, print func(a ...interface{})) error {
	timeout := time.Duration(conf.Timeout) * time.Second
	start := time.Now()
	address := conn.URL.String()
	print("Ping http address: " + address)
	if conf.Status > 0 {
		print("Expect HTTP status" + strconv.Itoa(conf.Status))
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	for k, v := range conf.Headers {
		print("Adding header " + k + ": " + v)
		req.Header.Add(k, v)
	}

	for {
		resp, err := client.Do(req)
		if resp != nil {
			print("Ping http address " + address + " " + resp.Status)
		}

		if err == nil {
			if conf.Status > 0 && conf.Status == resp.StatusCode {
				return nil
			} else if conf.Status == 0 && resp.StatusCode < http.StatusInternalServerError {
				return nil
			}
		}

		if time.Since(start) > timeout {
			return errors.New(resp.Status)
		}

		time.Sleep(time.Duration(conf.Retry) * time.Millisecond)
	}
}

// pingHost check if the host (hostname:port) is responding properly
func pingHost(conn *Connection, conf *Config, print func(a ...interface{})) error {
	timeout := time.Duration(conf.Timeout) * time.Second
	start := time.Now()
	address := conn.URL.Host
	print("Ping host: " + address)

	for {
		_, err := net.DialTimeout(conn.NetworkType, address, time.Second)
		print("Ping host: " + address)

		if err == nil {
			print("Up: " + address)
			return nil
		}

		print("Down: " + address)
		print(err)
		if time.Since(start) > timeout {
			return err
		}

		time.Sleep(time.Duration(conf.Retry) * time.Millisecond)
	}
}
