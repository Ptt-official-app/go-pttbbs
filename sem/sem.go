package sem

// from: https://github.com/shubhros/drunkendeluge/blob/master/semaphore/semaphore.go

// #include <sys/sem.h>
//
// /* https://comp.os.linux.development.system.narkive.com/rvJxp3Vb/union-variable-error-storage-size-isn-t-known */
// #if defined ( _SEM_SEMUN_UNDEFINED )
// union semun {
//   int val; /* value for SETVAL */
//   struct semid_ds *buf; /* buffer for IPC_STAT, IPC_SET */
//   unsigned short int *array; /* array for GETALL, SETALL */
//   struct seminfo *__buf; /* buffer for IPC_INFO */
// };
// #endif
//
//
// #ifndef SEM_R
// #define SEM_R 0400
// #endif
//
// #ifndef SEM_A
// #define SEM_A 0200
// #endif
// int semctlsetvalwrapper(int semid, int semnum, int val) {
//   union semun s;
//   s.val = val;
//   return semctl(semid, semnum, SETVAL, s);
// }
import "C"

import (
	"syscall"
	"unsafe"
)

const (
	IPC_CREAT = C.IPC_CREAT
	IPC_EXCL  = C.IPC_EXCL
	IPC_RMID  = C.IPC_RMID
	SETVAL    = C.SETVAL
	GETVAL    = C.GETVAL
	SEM_A     = C.SEM_A
	SEM_R     = C.SEM_R
	SEM_UNDO  = C.SEM_UNDO
)

type Semaphore struct {
	semid int
	nsems int
}

type semop struct {
	semNum  uint16
	semOp   int16
	semFlag int16
}

func errnoErr(errno syscall.Errno) error {
	switch errno {
	case syscall.Errno(0):
		return nil
	default:
		return errno
	}
}

func SemGet(key int, nsems int, flags int) (*Semaphore, error) {
	r1, _, errno := syscall.Syscall(syscall.SYS_SEMGET,
		uintptr(key), uintptr(nsems), uintptr(flags))
	if errno == syscall.Errno(0) {
		return &Semaphore{semid: int(r1), nsems: nsems}, nil
	} else {
		return nil, errnoErr(errno)
	}
}

func (s *Semaphore) Destroy() error {
	_, _, errno := syscall.Syscall(syscall.SYS_SEMCTL, uintptr(s.semid),
		uintptr(0), uintptr(IPC_RMID))
	return errnoErr(errno)
}

func (s *Semaphore) GetVal(semNum int) (int, error) {
	val, _, errno := syscall.Syscall(syscall.SYS_SEMCTL, uintptr(s.semid),
		uintptr(semNum), uintptr(GETVAL))
	return int(val), errnoErr(errno)
}

func (s *Semaphore) SetVal(semNum int, val int) error {
	cerrno := C.semctlsetvalwrapper(C.int(s.semid), C.int(semNum), C.int(val))
	return errnoErr(syscall.Errno(int(cerrno)))
}

func (s *Semaphore) Post(semNum int) error {
	post := semop{semNum: uint16(semNum), semOp: 1, semFlag: SEM_UNDO}
	_, _, errno := syscall.Syscall(syscall.SYS_SEMOP, uintptr(s.semid),
		uintptr(unsafe.Pointer(&post)), uintptr(s.nsems))
	return errnoErr(errno)

}

func (s *Semaphore) Wait(semNum int) error {
	wait := semop{semNum: uint16(semNum), semOp: -1, semFlag: SEM_UNDO}
	_, _, errno := syscall.Syscall(syscall.SYS_SEMOP, uintptr(s.semid),
		uintptr(unsafe.Pointer(&wait)), uintptr(s.nsems))
	return errnoErr(errno)
}
