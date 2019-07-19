package qio

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

type BufferType int

const (
	BufferMemory   BufferType = 1
	BufferTempFile BufferType = 2
)

type Buffer struct {
	bufferType     BufferType
	bufferFile     *os.File
	bufferMemory   *bytes.Buffer
	bufferFileName string
}

func NewBuffer(bufferType BufferType) *Buffer {
	var buffer = &Buffer{
		bufferType: bufferType,
	}

	if bufferType == BufferMemory {
		buffer.bufferMemory = &bytes.Buffer{}
	} else {
		var err error
		var dirs = []string{"temp", "tempdir", "tmp", "tmpdir", "tmpcache"}
		for _, dir := range dirs {
			buffer.bufferFile, err = ioutil.TempFile(dir, "buffer")
			if err != nil {
				err = os.Mkdir(dir, 0664)
				if err != nil {
					return nil
				}
				buffer.bufferFile, err = ioutil.TempFile("temp", "buffer")
			}
		}
		if buffer.bufferFile == nil {
			panic(fmt.Errorf("unable to create temp file in all such directorys %v , %v", dirs, err.Error()))
		}
		buffer.bufferFileName = buffer.bufferFile.Name()
	}
	return buffer
}

func (o *Buffer) Read(bytes []byte) (int, error) {
	if o.bufferType == BufferMemory {
		return o.bufferMemory.Read(bytes)
	} else {
		return o.bufferFile.Read(bytes)
	}
}

func (o *Buffer) Write(bytes []byte) (int, error) {
	if o.bufferType == BufferMemory {
		return o.bufferMemory.Write(bytes)
	} else {
		return o.bufferFile.Write(bytes)
	}
}

func (o *Buffer) WriteString(s string) (int, error) {
	if o.bufferType == BufferMemory {
		return o.bufferMemory.WriteString(s)
	} else {
		return o.bufferFile.WriteString(s)
	}
}

func (o *Buffer) Bytes() ([]byte, error) {
	if o.bufferType == BufferMemory {
		return o.bufferMemory.Bytes(), nil
	} else {
		var err = o.bufferFile.Sync()
		if err != nil {
			return nil, err
		}
		n, err := o.bufferFile.Seek(0, 2)
		if err == nil {

			_, err = o.bufferFile.Seek(0, 0)
			if err != nil {
				return err
			}
			var data = make([]byte, n)
			_, err := o.bufferFile.Read(data)
			return data, err
		} else {
			return nil, err
		}
	}
}

func (o *Buffer) Close() {
	if o.bufferType == BufferMemory {
		o.bufferMemory.Reset()
	} else {
		defer os.Remove(o.bufferFileName)
		_ = o.bufferFile.Close()
	}
}
