package main

import (
	"os"
)

type Datastore interface {
	Get(remote string) (*os.File, error)
	Put(local, remote string) error
	Delete(remote string) error
	Ls(remote string) []string
	Capacity() uint64
}
