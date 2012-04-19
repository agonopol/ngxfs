package main

import (
	"crypto/sha1"
	"encoding/binary"
	"hash"
	"math"
	"os"
)

type Server struct {
	weight uint
	addr   string
	store  DataStore
}

func (this *Server) Get(remote string) (*os.File, error) {
	return nil, nil
}

type Ring struct {
	servers []*Server
	crypt   hash.Hash
}

func NewRing(servers map[string]uint, store DataStore) *Ring {
	this := new(Ring)
	this.crypt = sha1.New()
	total := uint(0)
	for _, weight := range servers {
		total += weight
	}
	for addr, weight := range servers {
		times := math.Floor((float64(len(servers)) * float64(weight)) / float64(total))
		for i := 0; i < int(times); i++ {
			this.servers = append(this.servers, &Server{weight, addr, store})
		}

	}
	return this
}

func (this *Ring) Get(remote string) (*os.File, error) {
	key := this.crypt.Sum([]byte(remote))
	index := binary.BigEndian.Uint64(key)
	this.servers[index%uint64(len(this.servers))].Get(remote)
	return nil, nil
}

func (this *Ring) Put(local, remote string) error {
	return nil
}

func (this *Ring) Ls() []string {
	return []string{}
}
