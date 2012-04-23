package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var get *bool = flag.Bool("get", true, "get <remote>")
var put *bool = flag.Bool("put", false, "put <local> <remote>")
var del *bool = flag.Bool("del", false, "del <remote>")

func main() {
	flag.Parse()
	args := flag.Args()
	http := NewHttpDatastore("nutrun.com", 1)
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
	} else if *get {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		in, err := http.Get(args[0])
		defer in.Close()
		if err != nil {
			log.Fatal(err)
		}
		r := bufio.NewReader(in)
		out := bufio.NewWriter(os.Stdout)
		defer out.Flush()
		for {
			line, e := r.ReadSlice('\n')
			if e != nil {
				break
			}
			out.Write(line)
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
