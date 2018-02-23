package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type config struct {
	Auth struct {
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"auth"`
	Registry struct {
		URL      string `json:"url"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"registry"`
	Blacklist []string `json:"blacklist"`
}

func readConfig(file string) config {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Could not find the config file: %s\n", file)
	}
	var configFile config
	err = json.Unmarshal(data, &configFile)
	if err != nil {
		log.Fatalf("Could not parse config file: %v\n", err)
	}
	return configFile
}
