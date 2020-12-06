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
	flags := 0600 | IPC_CREAT | IPC_EXCL
	if isUseHugeTlb {
		flags |= SHM_HUGETLB
	}
	shmid, err = shmget(key, size, flags)
	log.Debugf("shm.CreateShm: after 1st shmget: shmid: %v err: (%v/%v)", shmid, err, reflect.TypeOf(err))

	isEExist := false
	if os.IsExist(err) {
		isEExist = true
		flags = 0600 | IPC_CREAT
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
	log.Debugf("shm.CloseShm: After detach shm: shmaddr: %v ret: %v err: %v", shmaddr, cret, err)

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

func ReadAt(shmaddr unsafe.Pointer, offset int, size uintptr, outptr unsafe.Pointer) {
	C.readwrapper(outptr, shmaddr, C.int(offset), C.ulong(size))
}

func WriteAt(shmaddr unsafe.Pointer, offset int, size uintptr, inptr unsafe.Pointer) {
	C.writewrapper(shmaddr, C.int(offset), inptr, C.ulong(size))
}

func IncUint32(shmaddr unsafe.Pointer, offset int) {
	C.incuint32wrapper(shmaddr, C.int(offset))
}

func SetOrUint32(shmaddr unsafe.Pointer, offset int, flag uint32) {
	C.set_or_uint32wrapper(shmaddr, C.int(offset), C.uint(flag))
}

func InnerSetInt32(shmaddr unsafe.Pointer, offsetSrc int, offsetDst int) {
	C.innerset_int32wrapper(shmaddr, C.int(offsetSrc), C.int(offsetDst))
}

func Memset(shmaddr unsafe.Pointer, offset int, c byte, size uintptr) {
	C.memsetwrapper(shmaddr, C.int(offset), C.uchar(c), C.ulong(size))
}

func SetBCACHEPTR(shmaddr unsafe.Pointer, offset int) {
	C.set_bcacheptr(shmaddr, C.int(offset))
}

func QsortCmpBoardName(shmaddr unsafe.Pointer, offset int, n uint32) {
	C.qsort_cmpboardname_wrapper(shmaddr, C.int(offset), C.ulong(n))
}

func QsortCmpBoardClass(shmaddr unsafe.Pointer, offset int, n uint32) {
	C.qsort_cmpboardclass_wrapper(shmaddr, C.int(offset), C.ulong(n))
}

//Memcmp
//
//memcmp(shmaddr+offset, cmpaddr, size)
//
//Return:
//	int: 0: shm == gomem, <0: shm < gomem, >0: shm > gomem
func Memcmp(shmaddr unsafe.Pointer, offset int, size uintptr, cmpaddr unsafe.Pointer) int {
	cret := C.memcmpwrapper(shmaddr, C.int(offset), C.ulong(size), cmpaddr)

	return int(cret)
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
