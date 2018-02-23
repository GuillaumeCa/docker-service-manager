package main

import (
	"flag"
)

var cfg config

func main() {
	updateFlag := flag.Bool("u", false, "update the services")
	flag.Parse()
	cfg = readConfig("config.json")

	if *updateFlag {
		m := newManager()
		m.updateServices(cfg)
	}

	newServer()
}
