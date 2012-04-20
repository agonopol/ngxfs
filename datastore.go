package main

import (
	"io"
)

type Datastore interface {
	Get(remote string) (io.ReadCloser, error)
	Put(local, remote string) error
	Delete(remote string) error
	Ls(remote string) []string
	Capacity() uint64
}
