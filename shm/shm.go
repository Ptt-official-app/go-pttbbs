package shm

//#include "shm.h"
import "C"

import (
	"os"
	"reflect"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

func CreateShm(key types.Key_t, size types.Size_t, isUseHugeTlb bool) (shmid int, shmaddr unsafe.Pointer, isNew bool, err error) {
	flags := 0o600 | IPC_CREAT | IPC_EXCL
	if isUseHugeTlb {
		flags |= SHM_HUGETLB
	}
	shmid, err = shmget(key, size, flags)
	log.Debugf("shm.CreateShm: after 1st shmget: shmid: %v err: (%v/%v)", shmid, err, reflect.TypeOf(err))

	isEExist := false
	if os.IsExist(err) {
		isEExist = true
		flags = 0o600 | IPC_CREAT
		if isUseHugeTlb {
			flags |= SHM_HUGETLB
		}
		shmid, err = shmget(key, size, flags)
		log.Debugf("shm.CreateShm: after 2nd shmget: shmid: %v err: (%v/%v)", shmid, err, reflect.TypeOf(err))
	}
	if shmid < 0 {
		log.Errorf("shm.CreateShm: unable to create shm: key: %v size: %v", key, size)
		return shmid, nil, false, err
	}

	shmaddr, err = shmat(shmid, nil, 0)
	log.Infof("shm.CreateShm: after shmat: shmaddr: %v e: %v", shmaddr, err)
	if err != nil {
		return -1, nil, false, err
	}

	return shmid, shmaddr, !isEExist, nil
}

func OpenShm(key types.Key_t, size types.Size_t, is_usehugetlb bool) (shmid int, shmaddr unsafe.Pointer, err error) {
	flags := 0
	if is_usehugetlb {
		flags |= SHM_HUGETLB
	}
	shmid, err = shmget(key, size, flags)

	if err != nil {
		log.Errorf("shm.OpenShm: unable to create shm: key: %v size: %v", key, size)
		return shmid, nil, err
	}

	shmaddr, err = shmat(shmid, nil, 0)
	if err != nil {
		return -1, nil, err
	}

	return shmid, shmaddr, nil
}

func CloseShm(shmid int, shmaddr unsafe.Pointer) (err error) {
	cret, err := C.shmdt(shmaddr)
	log.Infof("shm.CloseShm: After detach shm: shmaddr: %v ret: %v err: %v", shmaddr, cret, err)

	if err != nil {
		return err
	}

	cret, err = C.shmctl(C.int(shmid), C.IPC_RMID, nil)
	log.Infof("shm.CloseShm: After close shm: shmaddr: %v ret: %v, err: %v", shmaddr, cret, err)

	if int(cret) < 0 {
		return ErrUnableToCloseShm
	}

	return nil
}

func shmget(key types.Key_t, size types.Size_t, shmflg int) (int, error) {
	cshmid, err := C.shmget(C.int(key), C.ulong(size), C.int(shmflg))
	shmid := int(cshmid)
	if shmid < 0 {
		log.Errorf("unable to shmget: shmid: %v e: %v", shmid, err)
	}
	return shmid, err
}

func shmat(shmid int, shmaddr unsafe.Pointer, shmflg int) (unsafe.Pointer, error) {
	newShmAddr, err := C.shmat(C.int(shmid), shmaddr, C.int(shmflg))
	if err != nil {
		return nil, err
	}

	return newShmAddr, nil
}
