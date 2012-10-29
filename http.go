package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"path"
	"net/url"
	"log"
	"fmt"
    "strings"
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

func checkHttpResponse(resp *http.Response) error {
	if resp.StatusCode == 404 {
		return NotFoundError(resp.Status) 
	}
	if resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}
	return nil
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
	if err := checkHttpResponse(resp); err != nil {
		resp.Body.Close()
		return nil, err
	}
	return resp.Body, nil
}

func (this *HttpDatastore) Get(remote string) (io.ReadCloser, int64, error) {
	resp, err := http.Get(this.url(remote))
	if err != nil {
		return nil, 0, err
	}
	if err := checkHttpResponse(resp); err != nil {
		resp.Body.Close()
		return nil, 0, err
	}
	if resp.ContentLength == -1 {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("Unknown Content-Length when doing Get on url [%v]", remote)
	}
	return resp.Body, resp.ContentLength, nil
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
	if err := checkHttpResponse(resp); err != nil {
		resp.Body.Close()
		return nil, err
	}
	return resp.Body, nil
}

func (this *HttpDatastore) Ls(lspath string, recursive bool) ([]string, error) {
    if !strings.HasSuffix(lspath, "/") {
        lspath = lspath + "/"
    }
	resp, e := http.Get(this.url(lspath))
	if e != nil {
		return nil, e
	}
	if err := checkHttpResponse(resp); err != nil {
		resp.Body.Close()
		return nil, err
	}

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}
    resp.Body.Close()

	re, e := regexp.Compile(`\<a href\=\"(.*)\"\>`)
	if e != nil {
		panic(e)
	}
	results := re.FindAllSubmatch(body, -1)
	if len(results) == 0 {
		return make([]string, 0), nil
	}
	links := make([]string, 0)
	for _, result := range results {
		if file := string(result[1]); file != "../" {
			parsed, err := url.Parse(file)
			if err != nil {
				log.Panicf("Error parsing request uri [%v] file. err: %v", file, err)
			}

            if strings.HasSuffix(parsed.Path, "/") && recursive {
                subdir := path.Join(lspath, parsed.Path)
                sublinks, err := this.Ls(subdir, recursive)
                if err != nil {
                    return make([]string, 0), nil
                }
                for _, sublink := range(sublinks) {
                    links = append(links, path.Join(parsed.Path, sublink))
                }
            } else {
                links = append(links, parsed.Path)
            }
		}
	}

	return links, nil
}

func (this *HttpDatastore) Url(elem ...string) string {
	return "http://" + this.host + path.Join(elem...)
}

