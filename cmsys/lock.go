package cmsys

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
)

// pttLock
//
// XXX https://github.com/ptt/pttbbs/issues/100
// Need to sync with ptt-code.
// Requires file to be writable.
func pttLock(file *os.File, offset int64, theSize uintptr, mode int) (err error) {
	fd := file.Fd()

	lock_it := &syscall.Flock_t{
		Whence: int16(os.SEEK_CUR),
		Start:  offset,
		Len:    int64(theSize),
		Type:   int16(mode),
	}

	for idx := 0; idx < N_ITER_PTTLOCK; idx++ {
		err = syscall.FcntlFlock(fd, syscall.F_SETLKW, lock_it)
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(1 * time.Second))
	}

	return err
}

// GoPttLock
//
// Original PttLock has no effect with multi-thread process.
// We use single lock for now.
func GoPttLock(file *os.File, filename string, offset int64, theSize uintptr) (err error) {
	filenameOffset := fmt.Sprintf("%s%x", filename, offset)
	err = lockFD(filenameOffset)
	if err != nil {
		return err
	}

	return pttLock(file, offset, theSize, syscall.F_WRLCK)
}

// GoPttUnlock
//
// Original PttLock has no effect with multi-thread process.
// We use single lock for now.
func GoPttUnlock(file *os.File, filename string, offset int64, theSize uintptr) (err error) {
	filenameOffset := fmt.Sprintf("%s%x", filename, offset)
	defer func() { _ = unlockFD(filenameOffset) }()

	return pttLock(file, offset, theSize, syscall.F_UNLCK)
}

// GoFlock
//
// Original Flock has no effect with multi-thread process.
// We use single lock for now.
func GoFlock(fd uintptr, filename string) (err error) {
	err = lockFD(filename)
	if err != nil {
		return err
	}

	return syscall.Flock(int(fd), syscall.LOCK_EX)
}

// GoFlock
//
// Original Flock has no effect with multi-thread process.
// We use single lock for now.
func GoFlockExNb(fd uintptr, filename string) (err error) {
	err = lockFD(filename)
	if err != nil {
		return err
	}

	return syscall.Flock(int(fd), syscall.LOCK_EX|syscall.LOCK_NB)
}

// GoFunlock
//
// Original Flock has no effect with multi-thread process.
// We use single lock for now.
func GoFunlock(fd uintptr, filename string) (err error) {
	defer func() { _ = unlockFD(filename) }()

	return syscall.Flock(int(fd), syscall.LOCK_UN)
}

func lockFD(filenameOffset string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	_, ok := lockFDMap[filenameOffset]
	if ok {
		return ErrPttLock
	}

	theLock := &sync.Mutex{}

	theLock.Lock()
	lockFDMap[filenameOffset] = theLock

	return nil
}

func unlockFD(filenameOffset string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	theLock, ok := lockFDMap[filenameOffset]
	if !ok {
		return nil
	}

	theLock.Unlock()
	delete(lockFDMap, filenameOffset)

	return nil
}
