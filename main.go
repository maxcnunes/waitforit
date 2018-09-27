package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// VERSION is definded during the build
var VERSION string

// Config describes the connection config
type Config struct {
	Protocol string            `json:"proto"`
	Host     string            `json:"host"`
	Port     int               `json:"port"`
	Address  string            `json:"address"`
	Status   int               `json:"status"`
	Timeout  int               `json:"timeout"`
	Retry    int               `json:"retry"`
	Headers  map[string]string `json:"headers"`
}

// FileConfig describes the structure of the config json file
type FileConfig struct {
	Configs []Config
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	s := ""
	// for _, v := range i {
	// 	s += " " + v
	// }
	return s
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() { // nolint gocyclo
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\n  %s [options] [-- post-command]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "The options are:\n\n")
		flag.PrintDefaults()
	}

	var fheaders arrayFlags

	address := flag.String("address", "", "address (e.g. http://google.com or tcp://mysql_ip:mysql_port)")
	proto := flag.String("proto", "", "protocol to use during the connection")
	host := flag.String("host", "", "host to connect")
	port := flag.Int("port", 0, "port to connect")
	status := flag.Int("status", 0, "expected status that address should return (e.g. 200")
	timeout := flag.Int("timeout", 10, "seconds to wait until the address become available")
	retry := flag.Int("retry", 500, "milliseconds to wait between retries")
	printVersion := flag.Bool("v", false, "show the current version")
	debug := flag.Bool("debug", false, "enable debug")
	file := flag.String("file", "", "path of json file to read configs from")
	flag.Var(&fheaders, "header", "list of headers sent in the http(s) ping request")

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
		headers := make(map[string]string)
		if len(fheaders) > 0 {
			for _, v := range fheaders {
				result := strings.SplitN(v, ":", 2)
				if len(result) != 2 {
					continue
				}
				headers[result[0]] = strings.TrimLeft(result[1], " ")
			}
		}

		fc = FileConfig{
			Configs: []Config{
				{
					Protocol: *proto,
					Host:     *host,
					Port:     *port,
					Address:  *address,
					Status:   *status,
					Timeout:  *timeout,
					Retry:    *retry,
					Headers:  headers,
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
