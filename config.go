package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type config map[string]Datastore

func newConfig() config {
	configPath := os.Getenv("NGXFS_CONF")
	if configPath == "" {
		log.Fatalln("NGXFS_CONF undefined")
	}
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	conf := make(map[string]uint64)
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatal(err)
	}
	servers := make(config)
	for k, v := range conf {
		servers[k] = NewHttpDatastore(k, v)
	}
	return servers
}

var Config = newConfig()
