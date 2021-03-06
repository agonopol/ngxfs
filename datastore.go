package main

import (
	"io"
)

type Datastore interface {
	Get(remote string) (io.ReadCloser, int64, error)
	Put(local, remote string) (io.ReadCloser, error)
	Delete(remote string) (io.ReadCloser, error)
	Ls(path string, recursive bool) ([]string, error)
	Capacity() uint64
	Host() string
	Url(elem ...string) string
}
