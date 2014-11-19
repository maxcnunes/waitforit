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
		if err != nil {
			logDebug("Down...")

			elapsed := time.Since(start)
			if elapsed > timeout {
				return err
			}
		} else {
			logDebug("Up...")
			return nil
		}
	}

	return nil
}

func main() {
	host := flag.String("host", "localshot", "host to connect")
	port := flag.Int("port", 80, "port to connect")
	timeout := flag.Int("timeout", 10, "timeout to wait port be available")
	debug = flag.Bool("debug", false, "enable debug")

	flag.Parse()

	logDebug("starting")
	if err := Dial(*host, *port, *timeout); err != nil {
		log.Fatal(err)
	}
}
