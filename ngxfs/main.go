package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"ngxfs"
	"strings"
)

var get *bool = flag.Bool("get", false, "-get <remote>")
var put *bool = flag.Bool("put", false, "-put <local> <remote>")
var del *bool = flag.Bool("del", false, "-del <remote>")
var deldir *bool = flag.Bool("deldir", false, "-deldir <remote>")
var ls *bool = flag.Bool("ls", false, "-ls <path>")
var url *bool = flag.Bool("url", false, "-ls -url <path>")
var translate *bool = flag.Bool("translate", false, "-translate <path>")
var translateall *bool = flag.Bool("translateall", false, "-translateall <file>")

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
	} else if *deldir {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		body, err := ring.DeleteDir(args[0])
		if err != nil {
			log.Fatal(err)
			WriteBody(body, os.Stderr)
		}
	} else if *ls {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		results, err := ring.Ls(args[0], *url)
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
		
		outputFile := path.Base(args[0])
		if _, err := os.Lstat(outputFile); err == nil {
			log.Fatalf("Cannot perform get as outputfile [%v] already exists", outputFile)
		}

		body, err := ring.Get(args[0])
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}
		WriteBody(body, file)
	} else if *translate {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		for _, path := range ring.Translate(args[0]) {
			fmt.Println(path)
		}
	} else if *translateall {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		file, err := os.Open(args[0])
		defer file.Close()
		if err != nil {
			log.Panicf("Error opening file %v. err: %v", args[0], err)
		}
		buf := bufio.NewReader(file)
		var line string
		for line, err = buf.ReadString('\n'); err == nil; line, err = buf.ReadString('\n') {
			for _, path := range ring.Translate(strings.Trim(line, "\n")) {
				fmt.Println(path)
			}
		}
		if err != io.EOF || (err == io.EOF && line != "") {
			log.Panicf("Error parsing file. line: [%v], err: [%v]", line, err)
		}
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
