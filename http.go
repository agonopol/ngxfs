package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type HttpDatastore struct {
	host string
	cap  uint64
}

func NewHttpDatastore(host string, cap uint64) *HttpDatastore {
	return &HttpDatastore{host, cap}
}

func (this *HttpDatastore) Capacity() uint64 {
	return this.cap
}

func (this *HttpDatastore) url(path string) string {
	return fmt.Sprintf("http://%s%s", this.host, path)
}

func (this *HttpDatastore) Put(local, remote string) error {
	f, e := os.Open(local)
	if e != nil {
		return e
	}
	defer f.Close()
	r := bufio.NewReader(f)
	req, e := http.NewRequest("PUT", this.url(remote), r)
	if e != nil {
		return e
	}
	resp, e := http.DefaultTransport.RoundTrip(req)
	if e != nil {
		return e
	}
	log.Printf("%s,%s", resp.Status, resp.Body)
	return nil
}

func (this *HttpDatastore) Get(remote string) (io.ReadCloser, error) {
	resp, err := http.Get(this.url(remote))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (this *HttpDatastore) Delete(remote string) error {
	req, err := http.NewRequest("DELETE", this.url(remote), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	log.Printf("%s,%s", resp.Status, resp.Body)
	return err
}

func (this *HttpDatastore) Ls() []string {
	return make([]string, 0)
}