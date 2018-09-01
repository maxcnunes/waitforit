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

			ch <- DialConn(conn, conf.Timeout, conf.Retry, conf.Status, print)
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
func DialConn(conn *Connection, timeoutSeconds int, retryMseconds int, status int, print func(a ...interface{})) error {
	print("Waiting " + strconv.Itoa(timeoutSeconds) + " seconds")
	if err := pingHost(conn, timeoutSeconds, retryMseconds, print); err != nil {
		return err
	}

	if conn.URL.Scheme == "http" || conn.URL.Scheme == "https" {
		return pingAddress(conn, timeoutSeconds, retryMseconds, status, print)
	}

	return nil
}

// pingAddress check if the full address is responding properly
func pingAddress(conn *Connection, timeoutSeconds int, retryMseconds int, status int, print func(a ...interface{})) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	start := time.Now()
	address := conn.URL.String()
	print("Ping http address: " + address)
  if status > 0 {
		print("Expect HTTP status" + strconv.Itoa(status))
	}

	for {
		resp, err := http.Get(address)

		if resp != nil {
			print("Ping http address " + address + " " + resp.Status)
		}

		if err == nil {
			if status > 0 && status == resp.StatusCode {
				return nil
			} else if status == 0 && resp.StatusCode < http.StatusInternalServerError {
				return nil
			}
		}

		if time.Since(start) > timeout {
			return errors.New(resp.Status)
		}

		time.Sleep(time.Duration(retryMseconds) * time.Millisecond)
	}
}

// pingHost check if the host (hostname:port) is responding properly
func pingHost(conn *Connection, timeoutSeconds int, retryMseconds int, print func(a ...interface{})) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
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

		time.Sleep(time.Duration(retryMseconds) * time.Millisecond)
	}
}
