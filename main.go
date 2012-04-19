package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var put *string = flag.String("put", "", "Transfer <file> to remote site")
var del *bool = flag.Bool("del", false, "Delete from remote site")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	if *put != "" {
		f, e := os.Open(*put)
		if e != nil {
			log.Fatal(e)
		}
		defer f.Close()
		r := bufio.NewReader(f)
		req, e := http.NewRequest("PUT", args[0], r)
		if e != nil {
			log.Fatal(e)
		}
		resp, e := http.DefaultTransport.RoundTrip(req)
		if e != nil {
			log.Fatal(e)
		}
		fmt.Printf("%s", resp.Body)
	} else if *del {
		fmt.Printf("DEL %s\n", args[0])
	} else {
		fmt.Printf("GET %s\n", args[0])
	}
}
