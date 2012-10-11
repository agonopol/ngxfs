package ngxfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Configuration struct {
	Servers map[string]Datastore
	Redun   uint
}

func (this *Configuration) UnmarshalJSON(data []byte) error {
	this.Servers = make(map[string]Datastore)
	conf := make(map[string]interface{})
	err := json.Unmarshal(data, &conf)
	if err != nil {
		return err
	}
	for server, weight := range conf["hosts"].(map[string]interface{}) {
		this.Servers[server] = NewHttpDatastore(server, uint64(weight.(float64)))
	}
	if redun, found := conf["redun"]; found {
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

func newConfiguration() *Configuration {
	configurl := os.Getenv("NGXFS_CONF")
	if configurl == "" {
		log.Fatal("NGXFS_CONF undefined")
	}
	conf := &Configuration{make(map[string]Datastore), 1}
	resp, err := http.Get(configurl)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("%d status code when retrieving %s", resp.StatusCode, configurl)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}
