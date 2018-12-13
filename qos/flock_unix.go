// +build linux solaris darwin freebsd openbsd netbsd dragonfly

package qos

import (
	"os"
	"syscall"
)

// CheckLock tries to obtain an exclude lock on a lockfile and returns an error if one occurs
func (s *FileLock) Lock() error {

	// open/create lock file
	f, err := os.OpenFile(s.name, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	s.file = f
	// set the lock type to F_WRLCK, therefor the file has to be opened writable
	flock := syscall.Flock_t{
		Type: syscall.F_WRLCK,
		Pid:  int32(os.Getpid()),
	}
	// try to obtain an exclusive lock - FcntlFlock seems to be the portable *ix way
	if err := syscall.FcntlFlock(s.file.Fd(), syscall.F_SETLK, &flock); err != nil {
		return ErrAlreadyRunning
	}

	return nil
}

// TryUnlock unlocks, closes and removes the lockfile
func (s *FileLock) UnLock() error {
	if (s.file == nil) {
		return nil;
	}
	// set the lock type to F_UNLCK
	flock := syscall.Flock_t{
		Type: syscall.F_UNLCK,
		Pid:  int32(os.Getpid()),
	}
	if err := syscall.FcntlFlock(s.file.Fd(), syscall.F_SETLK, &flock); err != nil {
		return err;
	}
	if err := s.file.Close(); err != nil {
		return err;
	}
	if err := os.Remove(s.name); err != nil {
		return err;
	}
	return nil
}
