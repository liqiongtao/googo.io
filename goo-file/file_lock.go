package goo_file

import (
	"os"
	"syscall"
)

type FileLock struct {
	Filename string
	fh       *os.File
}

func (fl *FileLock) Lock() (err error) {
	if fl.Filename == "" {
		fl.Filename = ".lock"
	}

	fl.fh, err = os.Create(fl.Filename)
	if err != nil {
		return
	}

	err = syscall.Flock(int(fl.fh.Fd()), syscall.LOCK_EX)
	if err != nil {
		_ = fl.fh.Close()
		fl.fh = nil
	}

	return
}

func (fl *FileLock) UnLock() (err error) {
	if fl.fh == nil {
		return nil
	}
	defer fl.release()

	if err = syscall.Flock(int(fl.fh.Fd()), syscall.LOCK_UN); err != nil {
		return
	}

	return
}

func (fl *FileLock) release() {
	if fl.fh != nil {
		_ = fl.fh.Close()
		_ = os.Remove(fl.Filename)
		fl.fh = nil
	}
}
