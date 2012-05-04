package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"testing"
)

type DudStore struct {
	cap  uint64
	name string
}

func (this *DudStore) Get(remote string) (io.ReadCloser, error) {
	return nil, nil
}

func (this *DudStore) Put(local, remote string) (io.ReadCloser, error) {
	return nil, nil
}
func (this *DudStore) Ls(remote string) []string {
	return make([]string, 0)
}

func (this *DudStore) Capacity() uint64 {
	return this.cap
}

func (this *DudStore) Delete(remote string) (io.ReadCloser, error) {
	return nil, nil
}

func (this *DudStore) Host() string {
	return this.name
}

func TestRingCollisions(t *testing.T) {
	config := make(map[string]Datastore)
	config["A"] = &DudStore{10, "A"}
	config["B"] = &DudStore{2324, "B"}
	config["C"] = &DudStore{152, "C"}
	config["D"] = &DudStore{123, "D"}
	config["E"] = &DudStore{523, "E"}
	ring := NewRing(1, config)
	love := ring.server("/Circle/Of/Love")
	life := ring.server("/Circle/Of/Life")
	if love.(*DudStore).name != ring.server("/Circle/Of/Love").(*DudStore).name {
		t.Errorf("Failed to get the same server for /Circle/Of/Love")
	}
	if life.(*DudStore).name != ring.server("/Circle/Of/Life").(*DudStore).name {
		t.Errorf("Failed to get the same server for /Circle/Of/Life")

	}
}

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890 abcdefghijklmnopqrstuvwxyz~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"

func RandomNumber(max int) int {
	n := big.NewInt(int64(max))
	r, _ := rand.Int(rand.Reader, n)
	return int(r.Int64())
}

func RandomString(n int) string {
	str := make([]byte, n)
	for i := 0; i < n; i++ {
		str[i] = chars[RandomNumber(len(chars))]
	}
	return string(str)
}

func BenchmarkRingDistribution(t *testing.B) {
	config := make(map[string]Datastore)
	config["A"] = &DudStore{10, "A"}
	config["B"] = &DudStore{2324, "B"}
	config["C"] = &DudStore{152, "C"}
	config["D"] = &DudStore{123, "D"}
	config["E"] = &DudStore{523, "E"}
	distro := make(map[string]uint)
	ring := NewRing(1, config)
	for i := 0; i < 27610; i++ {
		url := RandomString(RandomNumber(260))
		server := ring.server(url)
		if _, found := distro[server.(*DudStore).name]; !found {
			distro[server.(*DudStore).name] = 0
		}
		distro[server.(*DudStore).name] += 1
	}
	fmt.Printf("Original distribution\n")
	for url, server := range config {
		fmt.Printf("%s - %v\n", url, float64(server.Capacity())/float64(2761))
	}
	fmt.Printf("File distribution\n")
	for url, server := range distro {
		fmt.Printf("%s - %v\n", url, float64(server)/float64(276100))
	}
}

func TestRingReduce(t *testing.T) {
	config := make(map[string]Datastore)
	config["A"] = &DudStore{10, "A"}
	config["B"] = &DudStore{2324, "B"}
	ring := NewRing(1, config)
	reduced := ring.ReduceRing(config["B"])
	if len(ring.servers) == len(reduced.servers) {
		t.Errorf("Reduced did not eliminate servers")
	}
	for _, server := range reduced.servers {
		if server.(*DudStore).name == config["B"].(*DudStore).name {
			t.Fatal("Reduce did not remove B from server list")
		}
	}
}
