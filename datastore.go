package main

import (
	"io"
)

type Datastore interface {
	Get(remote string) (io.ReadCloser, error)
	Put(local, remote string) (io.ReadCloser, error)
	Delete(remote string) (io.ReadCloser, error)
	Ls(remote string) []string
	Capacity() uint64
	Host() string
}
