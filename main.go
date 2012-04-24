package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var get *bool = flag.Bool("get", true, "get <remote>")
var put *bool = flag.Bool("put", false, "put <local> <remote>")
var del *bool = flag.Bool("del", false, "del <remote>")

func main() {
	flag.Parse()
	args := flag.Args()

	configPath := os.Getenv("NGXFS_CONF")

	if configPath == "" {
		log.Fatalln("NGXFS_CONF undefined")
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	config := make(map[string]uint64)
	err = json.Unmarshal(content, config)
	if err != nil {
		log.Fatal(err)
	}
	servers := make(map[string]Datastore)
	for k, v := range config {
		servers[k] = NewHttpDatastore(k, v)
	}
	ring := NewRing(servers)
	if *put {
		if len(args) != 2 {
			flag.Usage()
			os.Exit(1)
		}
		body, err := ring.Put(args[0], args[1])
		if err != nil {
			log.Fatal(err)
			WriteBody(body, os.Stderr)
		}
	} else if *del {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		body, err := ring.Delete(args[0])
		if err != nil {
			log.Fatal(err)
			WriteBody(body, os.Stderr)
		}
	} else if *get {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		body, err := ring.Get(args[0])
		if err != nil {
			log.Fatal(err)
		}
		WriteBody(body, os.Stdout)
	} else {
		flag.Usage()
		os.Exit(1)
	}
}

func WriteBody(in io.ReadCloser, out io.Writer) {
	defer in.Close()
	r := bufio.NewReader(in)
	bout := bufio.NewWriter(out)
	defer bout.Flush()
	for {
		line, e := r.ReadSlice('\n')
		if e != nil {
			break
		}
		bout.Write(line)
	}
}
