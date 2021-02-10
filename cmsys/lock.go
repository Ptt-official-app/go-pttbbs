package cmsys

import (
	"os"
	"syscall"
	"time"
)

//PttLock
//
//XXX https://github.com/ptt/pttbbs/issues/100
//Need to sync with ptt-code.
func PttLock(file *os.File, offset int64, theSize uintptr, mode int) (err error) {
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
