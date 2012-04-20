package main

import (
	"flag"
	"log"
	"os"
)

var put *bool = flag.Bool("put", false, "put <local> <remote>")
var del *bool = flag.Bool("del", false, "del <remote>")

func main() {
	flag.Parse()
	args := flag.Args()
	http := NewHttpDatastore("localhost:7777", 1)
	if *put {
		if len(args) != 2 {
			flag.Usage()
			os.Exit(1)
		}
		err := http.Put(args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else if *del {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		err := http.Delete(args[0])
		if err != nil {
			log.Fatal(err)
		}
	} else {

	}
}
