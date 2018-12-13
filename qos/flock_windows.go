// +build windows

package qos

import (
	"os"
)

// Filename returns an absolute filename, appropriate for the operating system


// CheckLock tries to obtain an exclude lock on a lockfile and returns an error if one occurs
func (s * FileLock) Lock() error {
	if (s.file != nil) {
		return nil;
	}
	var tmpname = s.name + ".pre";
	if err := os.Rename(s.name, tmpname); err != nil && !os.IsNotExist(err) {
		return err;
	}
	os.Rename(tmpname, s.name);
	file, err := os.OpenFile(s.name, os.O_EXCL|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	s.file = file
	return nil
}

// TryUnlock closes and removes the lockfile
func (s *FileLock) UnLock() error {
	if (s.file == nil) {
		return nil;
	}
	if err := s.file.Close(); err != nil {
		return err;
	}
	s.file = nil;
	var tmpname = s.name + ".pre";
	if err := os.Rename(s.name, tmpname); err != nil {
		return err;
	}
	os.Rename(tmpname, s.name);
	return nil
}
