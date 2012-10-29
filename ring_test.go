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

func (this *DudStore) Get(remote string) (io.ReadCloser, int64, error) {
	return nil, 0, nil
}

func (this *DudStore) Put(local, remote string) (io.ReadCloser, error) {
	return nil, nil
}

func (this *DudStore) Delete(remote string) (io.ReadCloser, error) {
	return nil, nil
}

func (this *DudStore) DeleteDir(remoteDir string) (io.ReadCloser, error) {
	return nil, nil
}

func (this *DudStore) Ls(path string, recursive bool) ([]string, error) {
	return nil, nil
}

func (this *DudStore) Capacity() uint64 {
	return this.cap
}

func (this *DudStore) Host() string {
	return this.name
}

func (this *DudStore) Url(elem ...string) string {
	return ""
}



func TestContinuumCollisions(t *testing.T) {
	config := make(map[string]Datastore)
	config["A"] = &DudStore{10, "A"}
	config["B"] = &DudStore{2324, "B"}
	config["C"] = &DudStore{152, "C"}
	config["D"] = &DudStore{123, "D"}
	config["E"] = &DudStore{523, "E"}
	ring := NewContinuum(config)
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

func TestContinuumGeneratesDifferentServers(t *testing.T) {
	config := make(map[string]Datastore)
	config["A"] = &DudStore{10, "A"}
	config["B"] = &DudStore{10, "B"}
	ring := NewContinuum(config)
	datastores := ring.RedundantServers("a", 1)
	if datastores[0].Host() != "A" {
		t.Fatalf("Expected A, Got %v", datastores[0].Host())
	}
	datastores = ring.RedundantServers("b", 1)
	if datastores[0].Host() != "B" {
		t.Fatalf("Expected B, Got %v", datastores[0].Host())
	}

}

func BenchmarkContinuumDistribution(t *testing.B) {
	config := make(map[string]Datastore)
	config["A"] = &DudStore{10, "A"}
	config["B"] = &DudStore{2324, "B"}
	config["C"] = &DudStore{152, "C"}
	config["D"] = &DudStore{123, "D"}
	config["E"] = &DudStore{523, "E"}
	distro := make(map[string]uint)
	ring := NewContinuum(config)
	for i := 0; i < 276100; i++ {
		url := RandomString(RandomNumber(260))
		server := ring.server(url)
		if _, found := distro[server.(*DudStore).name]; !found {
			distro[server.(*DudStore).name] = 0
		}
		distro[server.(*DudStore).name] += 1
	}
	fmt.Printf("\nOriginal distribution\n")
	for url, server := range config {
		fmt.Printf("%s - %v\n", url, float64(server.Capacity())/float64(2761))
	}
	fmt.Printf("File distribution\n")
	for url, server := range distro {
		fmt.Printf("%s - %v\n", url, float64(server)/float64(276100))
	}
}

func TestRContinuumReduce(t *testing.T) {
	config := make(map[string]Datastore)
	config["A"] = &DudStore{10, "A"}
	config["B"] = &DudStore{2324, "B"}
	ring := NewContinuum(config)
	reduced := ring.reduce(config["B"])
	if len(ring.config) == len(reduced.config) {
		t.Errorf("Reduced did not eliminate servers")
	}
	for _, datastore := range reduced.config {
		if datastore.(*DudStore).name == config["B"].(*DudStore).name {
			t.Fatal("Reduce did not remove B from server list")
		}
	}
}
