package ngxfs

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"math"
	"sort"
	"log"
	"strings"
)

const POINTS_PER_SERVER = 320

type multiReadCloser struct {
	readers io.Reader
	closers []io.ReadCloser
}

func MultiReadCloser(rclosers []io.ReadCloser) *multiReadCloser {
	readers := make([]io.Reader, len(rclosers))
	for i, reader := range rclosers {
		readers[i] = reader
	}
	return &multiReadCloser{io.MultiReader(readers...), rclosers}
}

func (this *multiReadCloser) Read(p []byte) (int, error) {
	return this.readers.Read(p)
}

func (this *multiReadCloser) Close() error {
	var err error
	for _, closer := range this.closers {
		e := closer.Close()
		if err != nil {
			err = e
		}
	}
	return err
}

type status struct {
	closer io.ReadCloser
	err    error
}

func newStatus(closer io.ReadCloser, err error) *status {
	return &status{closer, err}
}

type files struct {
	keys       []string
	err   	   error
	datastore  Datastore
}

func newfiles(keys []string, err error, datastore Datastore) *files {
	return &files{keys, err, datastore}
}

type ContinuumEntry struct {
	datastore Datastore
	hash      uint64
}

type ContinuumEntries []*ContinuumEntry

func (this ContinuumEntries) Len() int {
	return len(this)
}

func (this ContinuumEntries) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this ContinuumEntries) Less(i, j int) bool {
	return this[i].hash < this[j].hash
}

type Continuum struct {
	config  map[string]Datastore
	servers ContinuumEntries
	crypt   hash.Hash
}

func NewContinuum(servers map[string]Datastore) *Continuum {
	this := new(Continuum)
	this.config = servers
	this.crypt = sha1.New()
	return this
}

func (this *Continuum) hash(remote string) uint64 {
	this.crypt.Reset()
	if n, err := this.crypt.Write([]byte(remote)); n != len(remote) || err != nil {
		log.Panicf("Error writing to hash. string: %v. number of bytes written: %v. err: %v", remote, n, err)
	}
	return binary.BigEndian.Uint64(this.crypt.Sum(nil))
}

func (this *Continuum) build() {
	total := uint64(0)
	for _, server := range this.config {
		total += server.Capacity()
	}
	for _, server := range this.config {
		times := math.Floor((float64(len(this.config)) * POINTS_PER_SERVER * float64(server.Capacity())) / float64(total))
		for i := 0; i < int(times); i++ {
			hash := this.hash(fmt.Sprintf("%s:%d", server.Host(), i))
			this.servers = append(this.servers, &ContinuumEntry{server, hash})
		}

	}
	sort.Sort(this.servers)
}

func (this *Continuum) server(remote string) Datastore {
	index := this.hash(remote)
	if this.servers == nil { this.build() }
	return this.servers[index%uint64(len(this.servers))].datastore
}

func (this *Continuum) RedundantServers(remote string, redun uint) []Datastore {
	servers := make([]Datastore, redun)
	reduced := this
	for i := uint(0); i < redun; i++ {
		if i > 0 {
			reduced = reduced.reduce(servers[i-1])
		}
		servers[i] = reduced.server(remote)
	}
	return servers
}

func (this *Continuum) reduce(reducer Datastore) *Continuum {
	reduced := make(map[string]Datastore)
	for _, datastore := range this.config {
		if datastore.Host() != reducer.Host() {
			reduced[datastore.Host()] = datastore
		}
	}
	return NewContinuum(reduced)
}

type Ring struct {
	continuum *Continuum
	redun     uint
}

func NewRing(redun uint, servers map[string]Datastore) *Ring {
	return &Ring{NewContinuum(servers), redun}
}

func (this *Ring) Get(remote string) (io.ReadCloser, error) {
	var err error
	for _, server := range this.continuum.RedundantServers(remote, this.redun) {
		closer, e := server.Get(remote)
		err = e
		if err == nil {
			return closer, nil
		}
	}
	return nil, err
}

func (this *Ring) Delete(remote string) (io.ReadCloser, error) {
	var err error
	closers := make([]io.ReadCloser, this.redun)
	for i, server := range this.continuum.RedundantServers(remote, this.redun) {
		closers[i], err = server.Delete(remote)
		if err != nil {
			return MultiReadCloser(closers[0:i]), err
		}
	}
	return MultiReadCloser(closers), nil
}

func (this *Ring) Put(local, remote string) (io.ReadCloser, error) {
	closers := make([]io.ReadCloser, this.redun)
	stats := make(chan *status, this.redun)
	var err error
	for _, server := range this.continuum.RedundantServers(remote, this.redun) {
		go func(server Datastore) {
			stats <- newStatus(server.Put(local, remote))
		}(server)
	}
	for i := uint(0); i < this.redun; i++ {
		stat := <-stats
		closers[i] = stat.closer
		if stat.err != nil {
			err = stat.err
		}
	}
	close(stats)
	return MultiReadCloser(closers), err
}

func (this *Ring) Ls(path string, url bool) ([]string, error) {
	set := make(map[string]bool)
	links := make([]string, 0)
	lists := make(chan *files, len(this.continuum.config))
	for _, host := range this.continuum.config {
		go func(host Datastore) {
			keys, err := host.Ls(path)
			lists <- newfiles(keys, err, host)
		}(host)
	}
	for i := 0; i < len(this.continuum.config); i++ {
		results := <-lists
		if results.err != nil {
			return nil, results.err
		}
		for _, result := range results.keys {
			_, exists := set[result]
			if !exists {
				if url && !strings.HasSuffix(result, "/") {
					links = append(links, results.datastore.Url(path, result))
				} else {
					links = append(links, result)
				}
				set[result] = true
			}
		}
	}
	return links, nil
}

func (this *Ring) Translate(path string) []string {
	paths := make([]string, this.redun)
	for i, server := range this.continuum.RedundantServers(path, this.redun) {
		paths[i] = server.Url(path)
	}
	return paths
}
