package goo_utils

import (
	"os"
	"syscall"
)

type Lock interface {
	Lock() (err error)
	UnLock() (err error)
}

type FileLock struct {
	Filename string
	fh       *os.File
}

func (fl FileLock) Lock() (err error) {
	fl.fh, err = os.Create(fl.Filename)
	if err != nil {
		return
	}

	if err = syscall.Flock(int(fl.fh.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		fl.release()
		return
	}

	return
}

func (fl FileLock) UnLock() (err error) {
	defer fl.release()

	if err = syscall.Flock(int(fl.fh.Fd()), syscall.LOCK_UN); err != nil {
		return
	}

	return
}

func (fl FileLock) release() {
	if fl.fh != nil {
		fl.fh.Close()
	}
}
