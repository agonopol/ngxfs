package main

import (
	"os"
)

type DataStore interface {
	Get(remote string) (*os.File, error)
	Put(local, remote string) error
	Ls() []string
}
