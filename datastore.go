package main

import (
	"io"
)

type Datastore interface {
	Get(remote string) (io.ReadCloser, error)
	Put(local, remote string) (io.ReadCloser, error)
	Delete(remote string) (io.ReadCloser, error)
	Ls(path string) ([]string, error)
	Capacity() uint64
	Host() string
}
