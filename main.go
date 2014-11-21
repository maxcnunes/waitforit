package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var debug *bool

func logDebug(msg string) {
	if *debug {
		log.Println(msg)
	}
}

func Dial(host string, port int, timeoutSeconds int) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	start := time.Now()

	for {
		_, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err == nil {
			logDebug("Up...")
			return nil
		}

		logDebug("Down...")
		if time.Since(start) > timeout { return err }

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func main() {
	host := flag.String("host", "localhost", "host to connect")
	port := flag.Int("port", 80, "port to connect")
	timeout := flag.Int("timeout", 10, "time to wait until port become available")
	debug = flag.Bool("debug", false, "enable debug")

	flag.Parse()

	logDebug("Starting...")
	if err := Dial(*host, *port, *timeout); err != nil {
		log.Fatal(err)
	}
}
