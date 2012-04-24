package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
)

var get *bool = flag.Bool("get", true, "get <remote>")
var put *bool = flag.Bool("put", false, "put <local> <remote>")
var del *bool = flag.Bool("del", false, "del <remote>")

func main() {
	flag.Parse()
	args := flag.Args()
	http := NewHttpDatastore("0.0.0.0:8080", 1)
	if *put {
		if len(args) != 2 {
			flag.Usage()
			os.Exit(1)
		}
		body, err := http.Put(args[0], args[1])
		if err != nil {
			log.Fatal(err)
			WriteBody(body, os.Stderr)
		}
	} else if *del {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		body, err := http.Delete(args[0])
		if err != nil {
			log.Fatal(err)
			WriteBody(body, os.Stderr)
		}
	} else if *get {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		body, err := http.Get(args[0])
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
