package main

import (
	"crypto/sha1"
	"encoding/binary"
	goset "github.com/agonopol/goset"
	"hash"
	"io"
	"math"
)

type Ring struct {
	servers []Datastore
	crypt   hash.Hash
}

func NewRing(servers map[string]Datastore) *Ring {
	this := new(Ring)
	this.crypt = sha1.New()
	total := uint64(0)

	for _, server := range servers {
		total += server.Capacity()
	}

	for _, server := range servers {
		times := math.Floor((float64(len(servers)) * 320 * float64(server.Capacity())) / float64(total))
		for i := 0; i < int(times); i++ {
			this.servers = append(this.servers, server)
		}

	}
	return this
}

func (this *Ring) Get(remote string) (io.ReadCloser, error) {
	return this.server(remote).Get(remote)
}

func (this *Ring) Put(local, remote string) error {
	return this.server(remote).Put(local, remote)
}

func (this *Ring) Ls(remote string) []string {
	set := goset.New()
	for _, server := range this.servers {
		for _, path := range server.Ls(remote) {
			set.Add(path)
		}
	}
	paths := make([]string, 0)
	set.Do(func(path interface{}) {
		paths = append(paths, path.(string))
	})
	return paths
}

func (this *Ring) server(remote string) Datastore {
	key := this.crypt.Sum([]byte(remote))
	this.crypt.Reset()
	index := binary.BigEndian.Uint64(key)
	return this.servers[index%uint64(len(this.servers))]
}
