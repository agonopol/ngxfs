package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
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

func (this *HttpDatastore) Host() string {
	return this.host
}

func (this *HttpDatastore) Put(local, remote string) (io.ReadCloser, error) {
	f, e := os.Open(local)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	r := bufio.NewReader(f)
	req, e := http.NewRequest("PUT", this.url(remote), r)
	info, e := f.Stat()
	if e != nil {
		return nil, e
	}
	req.ContentLength = info.Size()
	if e != nil {
		return nil, e
	}
	resp, e := http.DefaultTransport.RoundTrip(req)
	if e != nil {
		return nil, e
	}
	if resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	return resp.Body, nil
}

func (this *HttpDatastore) Get(remote string) (io.ReadCloser, error) {
	resp, err := http.Get(this.url(remote))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	return resp.Body, nil
}

func (this *HttpDatastore) Delete(remote string) (io.ReadCloser, error) {
	req, err := http.NewRequest("DELETE", this.url(remote), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	return resp.Body, nil
}

func (this *HttpDatastore) Ls(path string) ([]string, error) {
	resp, e := http.Get(this.url(path))
	if e != nil {
		return nil, e
	}
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}
	re, e := regexp.Compile(`\<a href\=\"(.*)\"\>`)
	if e != nil {
		panic(e)
	}
	results := re.FindAllSubmatch(body, -1)
	links := make([]string, len(results))
	for i, result := range results {
		links[i] = string(result[1])
	}
	return links, nil
}
