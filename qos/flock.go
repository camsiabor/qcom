// package single provides a mechanism to ensure, that only one instance of a program is running

package qos

import (
	"github.com/camsiabor/qcom/util"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)


// Single represents the name and the open file descriptor
type FileLock struct {
	name string
	file * os.File
}

// New creates a Single instance
func NewFileLock(name string) *FileLock {
	return &FileLock{name: name}
}

func (o * FileLock) GetFileName() string {
	return o.name;
}

func (o * FileLock) GetFile() * os.File {
	return o.file;
}

func (o * FileLock) ReadString() (string, error) {
	var bytes, err = ioutil.ReadAll(o.file);
	if (err != nil) {
		return "", err;
	}
	return string(bytes[:]), err;
}

func (o * FileLock) WriteString(v interface{}) error {
	if (o.file == nil) {
		return errors.New("file not lock yet :" + o.name);
	}
	var s = util.AsStr(v, "");
	var _, err = o.file.WriteString(s);
	return err;
}
