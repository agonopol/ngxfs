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
	servers map[string]Datastore
	redun   uint
}

func (this *Configuration) UnmarshalJSON(data []byte) error {
	conf := make(map[string]interface{})
	err := json.Unmarshal(data, &conf)
	if err != nil {
		return err
	}
	for server, weight := range conf["hosts"].(map[string]interface{}) {
		this.servers[server] = NewHttpDatastore(server, uint64(weight.(float64)))
	}
	if redun, found := conf["redun"]; found {
		this.redun = uint(redun.(float64))
	}
	if len(this.servers) < int(this.redun) {
		return errors.New(fmt.Sprintf("Servers [%d] < Redun [%d]", len(this.servers), this.redun))
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
