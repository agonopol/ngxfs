package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Configuration struct {
	Servers map[string]Datastore
	Redun   uint
}

const HOSTS_KEY = "storageServers"
const REDUN_KEY = "redun"

func (this *Configuration) UnmarshalJSON(data []byte) error {
	this.Servers = make(map[string]Datastore)
	conf := make(map[string]interface{})
	err := json.Unmarshal(data, &conf)
	if err != nil {
		return err
	}
	hosts, found := conf[HOSTS_KEY]
	if !found {
		log.Panicf("No %v key found in configuration file", HOSTS_KEY)
	}
	for server, weight := range hosts.(map[string]interface{}) {
		server = strings.Replace(server, "http://", "", 1)
		this.Servers[server] = NewHttpDatastore(server, uint64(weight.(float64)))
	}
	if redun, found := conf[REDUN_KEY]; found {
		this.Redun = uint(redun.(float64))
	} else {
		this.Redun = 1
	}
	if this.Redun == 0 {
		return errors.New("0 Redun does not make sense")
	}
	if len(this.Servers) < int(this.Redun) {
		return errors.New(fmt.Sprintf("Servers [%d] < Redun [%d]", len(this.Servers), this.Redun))
	}
	return nil
}

func GetContent(url string) ([]byte, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != 200 {
			log.Fatalf("%d status code when retrieving %s", resp.StatusCode, url)
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}
	return ioutil.ReadFile(url)
}

func NewConfiguration() *Configuration {
	configurl := os.Getenv("NGXFS_CONF")
	if configurl == "" {
		log.Panic("NGXFS_CONF undefined")
	}

	content, err := GetContent(configurl)
	if err != nil {
		log.Fatal(err)
	}
	conf := &Configuration{make(map[string]Datastore), 1}
	if err := json.Unmarshal(content, &conf); err != nil {
		log.Fatal(err)
	}
	return conf
}
