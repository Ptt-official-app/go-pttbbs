package shm

//#include "shm.h"
import "C"
import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

func CreateShm(key types.Key_t, size types.Size_t, isUseHugeTlb bool) (shmid int, shmaddr unsafe.Pointer, isNew bool, err error) {
	flags := 0600 | IPC_CREAT | IPC_EXCL
	if isUseHugeTlb {
		flags |= SHM_HUGETLB
	}
	shmid = shmget(key, size, flags)

	isEExist := int(C.isEExist()) != 0
	if isEExist {
		flags = 0600 | IPC_CREAT
		if isUseHugeTlb {
			flags |= SHM_HUGETLB
		}
		shmid = shmget(key, size, flags)
	}
	if shmid < 0 {
		log.Errorf("shm.CreateShm: unable to create shm: key: %v size: %v", key, size)
		return shmid, nil, false, ErrInvalidShm
	}

	shmaddr, err = shmat(shmid, nil, 0)
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
	shmid = shmget(key, size, flags)

	if shmid < 0 {
		log.Errorf("shm.OpenShm: unable to create shm: key: %v size: %v", key, size)
		return shmid, nil, ErrInvalidShm
	}

	shmaddr, err = shmat(shmid, nil, 0)
	if err != nil {
		return -1, nil, err
	}

	return shmid, shmaddr, nil
}

func CloseShm(shmid int) (err error) {
	cerrno := C.shmctl(C.int(shmid), C.IPC_RMID, nil)

	log.Infof("After close shm: errno: %v", cerrno)

	if int(cerrno) < 0 {
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

func shmget(key types.Key_t, size types.Size_t, shmflg int) int {
	cshmid := C.shmget(C.int(key), C.ulong(size), C.int(shmflg))
	return int(cshmid)
}

func shmat(shmid int, shmaddr unsafe.Pointer, shmflg int) (unsafe.Pointer, error) {
	shmaddr = C.shmat(C.int(shmid), shmaddr, C.int(shmflg))
	if int(C.isPtrLessThan0(shmaddr)) != 0 {
		return nil, ErrUnableToAttachShm
	}

	return shmaddr, nil
}
