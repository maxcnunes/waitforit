package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

func dial(conn *Connection, timeoutSeconds int) error {
	logDebug("Waiting " + strconv.Itoa(timeoutSeconds) + " seconds")
	if err := pingTCP(conn, timeoutSeconds); err != nil {
		return err
	}

	if conn.Scheme != "http" && conn.Scheme != "https" {
		return nil
	}

	if err := pingHTTP(conn, timeoutSeconds); err != nil {
		return err
	}
	return nil
}

func parallelDial(conn *Connection, timeoutSeconds int, ch chan error) {
	err := dial(conn, timeoutSeconds)
	ch <- err
}

func pingHTTP(conn *Connection, timeoutSeconds int) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	start := time.Now()
	address := fmt.Sprintf("%s://%s:%d%s", conn.Scheme, conn.Host, conn.Port, conn.Path)
	logDebug("HTTP address: " + address)

	for {
		resp, err := http.Get(address)

		if resp != nil {
			logDebug("ping HTTP " + resp.Status)
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

func pingTCP(conn *Connection, timeoutSeconds int) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	start := time.Now()
	address := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
	logDebug("Dial address: " + address)

	for {
		_, err := net.DialTimeout(conn.Type, address, time.Second)
		logDebug("ping TCP")

		if err == nil {
			logDebug("Up")
			return nil
		}

		logDebug("Down")
		logDebug(err)
		if time.Since(start) > timeout {
			return err
		}

		time.Sleep(500 * time.Millisecond)
	}
}
