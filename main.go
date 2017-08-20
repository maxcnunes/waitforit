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
	Host     string `json:"host"`
	Port     int    `json:"port"`
	FullConn string `json:"fullConnection"`
	Timeout  int    `json:"timeout"`
}

// FileConfig describes the structure of the config json file
type FileConfig struct {
	Configs []Config
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\n  %s [options] [post-command]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "The options are:\n\n")
		flag.PrintDefaults()
	}

	fullConn := flag.String("full-connection", "", "full connection")
	host := flag.String("host", "", "host to connect")
	port := flag.Int("port", 80, "port to connect")
	timeout := flag.Int("timeout", 10, "time to wait until port become available")
	printVersion := flag.Bool("v", false, "show the current version")
	debug := flag.Bool("debug", false, "enable debug")
	file := flag.String("file", "", "path of json file to read configs from")

	flag.Parse()

	if *printVersion {
		fmt.Println("waitforit version " + VERSION)
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
					Host:     *host,
					Port:     *port,
					FullConn: *fullConn,
					Timeout:  *timeout,
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
	args := flag.Args()

	if len(args) == 0 {
		return nil
	}

	cmd := exec.Command(args[0], args[1:len(args)]...)
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
