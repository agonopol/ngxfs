package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var outputFile *string = flag.String("o", "", "-o <outputfile> <url>")
var put *bool = flag.Bool("put", false, "-put <local> <remote>")
var del *bool = flag.Bool("del", false, "-del <remote>")
var ls *bool = flag.Bool("ls", false, "-ls <path>")
var fullurl *bool = flag.Bool("url", false, "-url -ls <path>")
var translate *bool = flag.Bool("translate", false, "-translate <path>")
var translateall *bool = flag.Bool("translateall", false, "-translateall <file>")

func main() {
	// Remove all info from log output
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()
	config := NewConfiguration()
	ring := NewRing(config.Redun, config.Servers)
	if *put {
		if len(args) != 2 {
			flag.Usage()
			os.Exit(1)
		}
		body, err := ring.Put(args[0], args[1])
		if err != nil {
			WriteBody(body, os.Stderr)
			log.Fatal(err)
		}
		body.Close()
	} else if *del {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		body, err := ring.Delete(args[0])
		if err != nil {
			WriteBody(body, os.Stderr)
			log.Fatal(err)
		}
		body.Close()
	} else if *ls {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		results, err := ring.Ls(args[0], *fullurl)
		if err != nil {
			log.Fatal(err)
		}
		for _, result := range results {
			fmt.Println(result)
		}
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
		// Do a Get
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		
		var body io.ReadCloser
		var size int64
		var out io.WriteCloser
		var err error

		body, size, err = ring.Get(args[0])
		if err != nil {
			log.Fatal(err)
		}

		if *outputFile == "" {
			out = os.Stdout		
		} else {
			out, err = os.Create(*outputFile)
			defer out.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
		if bytesWritten := WriteBody(body, out); bytesWritten != size {
			log.Panicf("Bytes written [%v] does not equal expected size [%v]", bytesWritten, size)
		}
	} 
}

func WriteBody(in io.ReadCloser, out io.Writer) int64 {
	defer in.Close()
	n, err := io.Copy(out, in)
	if err != nil {
		log.Panicf("Error copying. err: %v", err)
	}
	return n
}
