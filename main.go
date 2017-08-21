package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// VERSION is definded during the build
var VERSION string

// Config describes the connection config
type Config struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Address string `json:"address"`
	Timeout int    `json:"timeout"`
}

// FileConfig describes the structure of the config json file
type FileConfig struct {
	Configs []Config
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\n  %s [options] [-- post-command]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "The options are:\n\n")
		flag.PrintDefaults()
	}

	address := flag.String("address", "", "address (e.g. http://google.com or tcp://mysql_ip:mysql_port)")
	host := flag.String("host", "", "host to connect")
	port := flag.Int("port", 80, "port to connect")
	timeout := flag.Int("timeout", 10, "time to wait until the address become available")
	printVersion := flag.Bool("v", false, "show the current version")
	debug := flag.Bool("debug", false, "enable debug")
	file := flag.String("file", "", "path of json file to read configs from")

	flag.Parse()

	if *printVersion {
		fmt.Println("waitforit version " + VERSION)
		return
	}

	if *address == "" && *host == "" && *file == "" {
		fmt.Fprintf(os.Stderr, "Missing : address, host or file field\n")
		flag.Usage()
		return
	}

	print := func(a ...interface{}) {}
	if *debug {
		print = func(a ...interface{}) {
			log.Print(a...)
		}
	}

	var fc FileConfig
	if *file != "" {
		if err := loadFileConfig(*file, &fc); err != nil {
			log.Fatal(err)
		}
	} else {
		fc = FileConfig{
			Configs: []Config{
				{
					Host:    *host,
					Port:    *port,
					Address: *address,
					Timeout: *timeout,
				},
			},
		}
	}

	if err := DialConfigs(fc.Configs, print); err != nil {
		log.Fatal(err)
	}

	if err := runPostCommand(); err != nil {
		os.Exit(1)
	}
}

func runPostCommand() error {
	extraArgs := flag.Args()
	nExtraArgs := len(extraArgs)
	if nExtraArgs == 0 {
		return nil
	}

	// Ensure with explict argument "--" is enabling a post command
	allArgs := os.Args
	nAllArgs := len(allArgs)
	if allArgs[nAllArgs-(nExtraArgs+1)] != "--" {
		return nil
	}

	cmd := exec.Command(extraArgs[0], extraArgs[1:len(extraArgs)]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func loadFileConfig(path string, fc *FileConfig) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(f).Decode(&fc); err != nil {
		return err
	}

	return nil
}
