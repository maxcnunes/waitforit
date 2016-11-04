package main

import (
	"flag"
	"fmt"
	"log"
)

// VERSION is definded during the build
var VERSION string
var debug *bool

func logDebug(msg interface{}) {
	if *debug {
		log.Print(msg)
	}
}

func main() {
	fullConn := flag.String("full-connection", "", "full connection")
	host := flag.String("host", "", "host to connect")
	port := flag.Int("port", 80, "port to connect")
	timeout := flag.Int("timeout", 10, "time to wait until port become available")
	printVersion := flag.Bool("v", false, "show the current version")
	debug = flag.Bool("debug", false, "enable debug")
	file := flag.String("file", "", "path of json file to read configs from")

	flag.Parse()

	if *printVersion {
		fmt.Println("waitforit version " + VERSION)
		return
	}

	if *file != "" {
		if err := useFileConfig(*file); err != nil {
			log.Fatal(err)
		}
		return
	}

	conn := buildConn(*host, *port, *fullConn)
	if conn == nil {
		log.Fatal("Invalid connection")
	}

	if err := dial(conn, *timeout); err != nil {
		log.Fatal(err)
	}
}
