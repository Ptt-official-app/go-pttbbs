package cmsys

import (
	"os"
	"sync"
	"syscall"
	"time"
)

//pttLock
//
//XXX https://github.com/ptt/pttbbs/issues/100
//Need to sync with ptt-code.
//Requires file to be writable.
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

//GoPttLock
//
//Original PttLock has no effect with multi-thread process.
//We use single lock for now.
func GoPttLock(file *os.File, filename string, offset int64, theSize uintptr) (err error) {

	lockFD(file.Fd())

	return pttLock(file, offset, theSize, syscall.F_WRLCK)
}

//GoPttUnlock
//
//Original PttLock has no effect with multi-thread process.
//We use single lock for now.
func GoPttUnlock(file *os.File, filename string, offset int64, theSize uintptr) (err error) {

	defer unlockFD(file.Fd())

	return pttLock(file, offset, theSize, syscall.F_UNLCK)
}

//GoFlock
//
//Original Flock has no effect with multi-thread process.
//We use single lock for now.
func GoFlock(fd uintptr, filename string) (err error) {
	lockFD(fd)

	return syscall.Flock(int(fd), syscall.LOCK_EX)
}

//GoFlock
//
//Original Flock has no effect with multi-thread process.
//We use single lock for now.
func GoFlockExNb(fd uintptr, filename string) (err error) {
	lockFD(fd)

	return syscall.Flock(int(fd), syscall.LOCK_EX|syscall.LOCK_NB)
}

//GoFunlock
//
//Original Flock has no effect with multi-thread process.
//We use single lock for now.
func GoFunlock(fd uintptr, filename string) (err error) {
	defer unlockFD(fd)

	return syscall.Flock(int(fd), syscall.LOCK_UN)
}

func lockFD(fd uintptr) (err error) {
	lock.Lock()
	defer lock.Unlock()

	_, ok := lockFDMap[fd]
	if ok {
		return ErrPttLock
	}

	theLock := &sync.Mutex{}

	theLock.Lock()
	lockFDMap[fd] = theLock

	return nil
}

func unlockFD(fd uintptr) (err error) {
	lock.Lock()
	defer lock.Unlock()

	theLock, ok := lockFDMap[fd]
	if !ok {
		return ErrPttLock
	}

	theLock.Unlock()
	delete(lockFDMap, fd)

	return nil
}
