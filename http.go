package ngxfs

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"path"
	"strings"
	"net/url"
	"log"
	"strconv"
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
	url := url.URL{Scheme: "http", Host: this.host, Path: path}
	return url.String()
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

func (this *HttpDatastore) Get(remote string) (io.ReadCloser, int64, error) {
	resp, err := http.Get(this.url(remote))
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode >= 300 {
		return nil, 0, errors.New(resp.Status)
	}
	
	contentLengthStr := resp.Header.Get("Content-Length")
	if contentLengthStr == "" {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("Did not find Content-Length when trying to Get [%v]", remote)
	}
	size, err := strconv.ParseInt(contentLengthStr, 10, 64)
	if err != nil {
		resp.Body.Close()
		return nil, 0, err
	}
	return resp.Body, size, nil
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

func (this *HttpDatastore) DeleteDir(remoteDir string) (io.ReadCloser, error) {
	if !strings.HasSuffix(remoteDir, "/") {
		return nil, fmt.Errorf("remoteDir [%v] does not end with a /", remoteDir)
	}

	req, err := http.NewRequest("DELETE", this.url(remoteDir), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 404 && resp.StatusCode >= 300 {
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
	if len(results) == 0 {
		return make([]string, 0), nil
	}
	links := make([]string, len(results) - 1)
	idx := 0
	for _, result := range results {
		if file := string(result[1]); file != "../" {
			parsed, err := url.Parse(file)
			if err != nil {
				log.Panicf("Error parsing request uri [%v] file. err: %v", file, err)
			}
			links[idx] = parsed.Path
			idx += 1
		}
	}
	return links, nil
}

func (this *HttpDatastore) Url(elem ...string) string {
	return "http://" + this.host + path.Join(elem...)
}

