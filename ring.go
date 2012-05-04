package main

import (
	"crypto/sha1"
	"encoding/binary"
	"hash"
	"io"
	"math"
	"sort"
)

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

type Ring struct {
	config    map[string]Datastore
	servers   []Datastore
	crypt     hash.Hash
	redundant uint
}

func (this *Ring) sortMapKeys(in map[string]Datastore) []string {
	keys := make([]string, 0)
	for k, _ := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func NewRing(redundancy uint, servers map[string]Datastore) *Ring {
	this := new(Ring)
	this.config = servers
	this.redundant = redundancy
	this.crypt = sha1.New()
	total := uint64(0)
	for _, server := range servers {
		total += server.Capacity()
	}
	for _, server := range this.sortMapKeys(servers) {
		times := math.Floor((float64(len(servers)) * 320 * float64(servers[server].Capacity())) / float64(total))
		for i := 0; i < int(times); i++ {
			this.servers = append(this.servers, servers[server])
		}

	}
	return this
}

func (this *Ring) ReduceRing(reducer Datastore) *Ring {
	reduced := make(map[string]Datastore)
	for _, server := range this.servers {
		if server.Host() != reducer.Host() {
			reduced[server.Host()] = server
		}
	}
	return NewRing(this.redundant-1, reduced)
}

func (this *Ring) Get(remote string) (io.ReadCloser, error) {
	var err error
	for _, server := range this.redudantServers(remote) {
		closer, err := server.Get(remote)
		if err != nil {
			return closer, nil
		}
	}
	return nil, err
}

func (this *Ring) Delete(remote string) (io.ReadCloser, error) {
	var err error
	closers := make([]io.ReadCloser, this.redundant)
	for i, server := range this.redudantServers(remote) {
		closers[i], err = server.Delete(remote)
		if err != nil {
			return MultiReadCloser(closers[0:i]), err
		}
	}
	return MultiReadCloser(closers), nil

}

func (this *Ring) Put(local, remote string) (io.ReadCloser, error) {
	var err error
	closers := make([]io.ReadCloser, this.redundant)
	for i, server := range this.redudantServers(remote) {
		closers[i], err = server.Put(local, remote)
		if err != nil {
			return MultiReadCloser(closers[0:i]), err
		}
	}
	return MultiReadCloser(closers), nil
}

func (this *Ring) Ls(remote string) []string {
	set := make(map[string]bool)
	for _, server := range this.config {
		for _, path := range server.Ls(remote) {
			set[path] = true
		}
	}
	paths := make([]string, 0)
	for path, _ := range set {
		paths = append(paths, path)
	}
	return paths
}

func (this *Ring) redudantServers(remote string) []Datastore {
	servers := make([]Datastore, this.redundant)
	reduced := this
	for i := uint(0); i < uint(this.redundant); i++ {
		servers[i] = reduced.server(remote)
		reduced = reduced.ReduceRing(servers[i])
	}
	return servers
}

func (this *Ring) server(remote string) Datastore {
	key := this.crypt.Sum([]byte(remote))
	this.crypt.Reset()
	index := binary.BigEndian.Uint64(key)
	return this.servers[index%uint64(len(this.servers))]
}
