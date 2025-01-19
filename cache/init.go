package cache

import (
	"time"

	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

// Init
// This is to init SHM with Version and Size checked.
func Init(key types.Key_t, isUseHugeTlb bool, isCreate bool) (err error) {
	err = InitMAP()
	if err != nil {
		return err
	}

	if types.IS_ALL_GUEST {
		return nil
	}

	return InitSHM(key, isUseHugeTlb, isCreate)
}

func InitMAP() (err error) {
	if MAP != nil {
		return nil
	}

	MAP, err = NewMap()

	return err
}

func InitSHM(key types.Key_t, isUseHugeTlb bool, isCreate bool) (err error) {
	if SHM != nil {
		return nil
	}

	SHM, err = NewShm(key, isUseHugeTlb, isCreate)

	return err
}

// Close
//
// XXX [WARNING] know what you are doing before using Close!.
// This is to be able to close the shared mem for the completeness of the mem-usage.
// However, in production, we create shm without the need of closing the shm.
//
// We simply use ipcrm to delete the shm if necessary.
//
// Currently used only in test.
// XXX not doing close shm to prevent opening too many shms in tests.
func CloseSHM() error {
	if !IsTest {
		return ErrInvalidOp
	}

	if SHM == nil {
		// Already Closed
		log.Errorf("cache.CloseSHM: already closed")
		return ErrShmNotInit
	}

	err := SHM.Close()
	if err != nil {
		log.Errorf("cache.CloseSHM: unable to close: e: %v", err)
		return err
	}

	SHM = nil

	time.Sleep(3 * time.Millisecond)

	log.Infof("cache.CloseSHM: done")

	return nil
}

func CloseMAP() (err error) {
	if !IsTest {
		return ErrInvalidOp
	}
	if MAP == nil {
		return ErrMapNotInit
	}

	MAP = nil

	return nil
}
