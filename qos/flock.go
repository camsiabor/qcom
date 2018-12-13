// package single provides a mechanism to ensure, that only one instance of a program is running

package qos

import (
	"os"
)


// Single represents the name and the open file descriptor
type FileLock struct {
	name string
	file * os.File
}

// New creates a Single instance
func New(name string) *FileLock {
	return &FileLock{name: name}
}
