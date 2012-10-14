package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"ngxfs"
	"strings"
)

var get *bool = flag.Bool("get", true, "get <remote>")
var put *bool = flag.Bool("put", false, "put <local> <remote>")
var del *bool = flag.Bool("del", false, "del <remote>")
var ls *bool = flag.Bool("ls", false, "ls <path>")
var translate *bool = flag.Bool("translate", false, "translate")

func next_put_pair(buf *bufio.Reader) (string, string, error) {
	local, err:= buf.ReadString('\n') 
	if err != nil {
		return local, "", err
	}
	remote, err := buf.ReadString('\n')
	if err == io.EOF {
		return local, remote, fmt.Errorf("Got EOF while trying to read the remote url. err: %v", err)
	}
	return strings.Trim(local, "\n"), strings.Trim(remote, "\n"), err
}

func main() {
	// Remove all info from log output
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()
	config := ngxfs.NewConfiguration()
	ring := ngxfs.NewRing(config.Redun, config.Servers)
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
	} else if *ls {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		results, err := ring.Ls(args[0])
		if err != nil {
			log.Fatal(err)
		}
		for _, result := range results {
			fmt.Println(result)
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
