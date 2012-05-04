package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Configuration struct {
	servers    map[string]Datastore
	redundancy uint
}

func (this *Configuration) UnmarshalJSON(data []byte) error {
	conf := make(map[string]interface{})
	err := json.Unmarshal(data, &conf)
	if err != nil {
		return err
	}
	for server, weight := range conf["servers"].(map[string]interface{}) {
		this.servers[server] = NewHttpDatastore(server, uint64(weight.(float64)))
	}
	if redundancy, found := conf["redundancy"]; found {
		this.redundancy = uint(redundancy.(float64))
	}
	if len(this.servers) < int(this.redundancy) {
		return errors.New(fmt.Sprintf("Servers [%d] < Redundancy [%d]", len(this.servers), this.redundancy))
	}
	return nil
}

func newConfiguration() *Configuration {
	configPath := os.Getenv("NGXFS_CONF")
	conf := &Configuration{make(map[string]Datastore), 1}
	if configPath == "" {
		log.Printf("NGXFS_CONF undefined")
	} else {
		content, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(content, &conf)

		if err != nil {
			log.Fatal(err)
		}
	}
	return conf
}

var Config = newConfiguration()
