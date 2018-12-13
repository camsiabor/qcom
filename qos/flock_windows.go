// +build windows

package qos

import (
	"os"
)

// Filename returns an absolute filename, appropriate for the operating system


// CheckLock tries to obtain an exclude lock on a lockfile and returns an error if one occurs
func (o * FileLock) Lock() error {
	if (o.file != nil) {
		return nil;
	}
	var tmpname = o.name + ".pre";
	if err := os.Rename(o.name, tmpname); err != nil && !os.IsNotExist(err) {
		return err;
	}
	os.Rename(tmpname, o.name);
	file, err := os.OpenFile(o.name, os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	o.file = file
	return nil
}

// TryUnlock closes and removes the lockfile
func (o *FileLock) UnLock() error {
	if (o.file == nil) {
		return nil;
	}
	if err := o.file.Close(); err != nil {
		return err;
	}
	o.file = nil;
	var tmpname = o.name + ".pre";
	if err := os.Rename(o.name, tmpname); err != nil {
		return err;
	}
	os.Rename(tmpname, o.name);
	return nil
}
