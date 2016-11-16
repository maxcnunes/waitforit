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
			conn := BuildConn(conf.Host, conf.Port, conf.FullConn)
			if conn == nil {
				ch <- fmt.Errorf("Invalid connection %#v", conf)
				return
			}

			ch <- DialConn(conn, conf.Timeout, print)
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
func DialConn(conn *Connection, timeoutSeconds int, print func(a ...interface{})) error {
	print("Waiting " + strconv.Itoa(timeoutSeconds) + " seconds")
	if err := pingTCP(conn, timeoutSeconds, print); err != nil {
		return err
	}

	if conn.Scheme != "http" && conn.Scheme != "https" {
		return nil
	}

	return pingHTTP(conn, timeoutSeconds, print)
}

func pingHTTP(conn *Connection, timeoutSeconds int, print func(a ...interface{})) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	start := time.Now()
	address := fmt.Sprintf("%s://%s:%d%s", conn.Scheme, conn.Host, conn.Port, conn.Path)
	print("HTTP address: " + address)

	for {
		resp, err := http.Get(address)

		if resp != nil {
			print("ping HTTP " + address + " " + resp.Status)
		}

		if err == nil && resp.StatusCode < http.StatusInternalServerError {
			return nil
		}

		if time.Since(start) > timeout {
			return errors.New(resp.Status)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func pingTCP(conn *Connection, timeoutSeconds int, print func(a ...interface{})) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	start := time.Now()
	address := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
	print("Dial address: " + address)

	for {
		_, err := net.DialTimeout(conn.Type, address, time.Second)
		print("ping TCP: " + address)

		if err == nil {
			print("Up: " + address)
			return nil
		}

		print("Down: " + address)
		print(err)
		if time.Since(start) > timeout {
			return err
		}

		time.Sleep(500 * time.Millisecond)
	}
}
