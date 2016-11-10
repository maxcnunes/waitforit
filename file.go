package main

import (
	"encoding/json"
	"os"
)

type singleConfig struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	ConnectionString string `json:"connectionString"`
	Timeout          int    `json:"timeout"`
}

// fileConfigs describes the structure of the config json file
type fileConfig struct {
	Configs []singleConfig
}

// handleFileConfig is the handler function for the file options
func handleFileConfig(config fileConfig) error {
	ch := make(chan error)
	for _, conn := range config.Configs {
		connection := buildConn(conn.Host, conn.Port, conn.ConnectionString)
		go parallelDial(connection, conn.Timeout, ch)
	}
	for i := 0; i < len(config.Configs); i++ {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}

// useFileConfig is the main function used to read and proccess the json config
func useFileConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	var config fileConfig
	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(&config); err != nil {
		return err
	}
	if err := handleFileConfig(config); err != nil {
		return err
	}
	return nil
}
