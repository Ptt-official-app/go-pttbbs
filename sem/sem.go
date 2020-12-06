package sem

// from: https://github.com/shubhros/drunkendeluge/blob/master/semaphore/semaphore.go

// #include "sem.h"
import "C"

import (
	log "github.com/sirupsen/logrus"
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

func SemGet(key int, nsems int, flags int) (*Semaphore, error) {
	cret, err := C.semget(C.int(key), C.int(nsems), C.int(flags))
	log.Debugf("sem.SemGet: key: %v nsems: %v cret: %v e: %v", key, nsems, cret, err)
	if err != nil {
		return nil, err
	}

	return &Semaphore{semid: int(cret), nsems: nsems}, nil
}

func (s *Semaphore) Destroy(semNum int) error {
	cret, err := C.semdestroywrapper(C.int(s.semid), C.int(semNum))

	log.Debugf("sem.Destroy: semid: %v cret: %v e: %v", s.semid, cret, err)
	return err
}

func (s *Semaphore) GetVal(semNum int) (int, error) {
	cret, err := C.semgetvalwrapper(C.int(s.semid), C.int(semNum))
	log.Debugf("sem.GetVal: semid: %v semNum: %v cret: %v e: %v", s.semid, semNum, cret, err)
	return int(cret), err
}

func (s *Semaphore) SetVal(semNum int, val int) error {
	cret, err := C.semctlsetvalwrapper(C.int(s.semid), C.int(semNum), C.int(val))
	log.Debugf("sem.SetVal: semid: %v semNum: %v val: %v cret: %v e: %v", s.semid, semNum, val, cret, err)
	return err
}

func (s *Semaphore) Post(semNum int) error {
	cret, err := C.sempostwrapper(C.int(s.semid), C.int(semNum))
	log.Debugf("sem.Post: semid: %v semNum: %v cret: %v e: %v", s.semid, semNum, cret, err)

	return err
}

func (s *Semaphore) Wait(semNum int) error {
	cret, err := C.semwaitwrapper(C.int(s.semid), C.int(semNum))
	log.Debugf("sem.Wait: semid: %v semNum: %v cret: %v err: %v", s.semid, semNum, cret, err)

	return err
}
