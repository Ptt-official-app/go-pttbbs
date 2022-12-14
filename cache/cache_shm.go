package cache

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

// AttachSHM
//
// This is to have attach_shm (no checking loaded), not attach_SHM (checking loaded)
// Should be used after LoadUHash (shmctl init) is done.
// Should be used only once in the beginning of the program.
func AttachSHM() error {
	if Shm != nil {
		return nil
	}

	return NewSHM(types.Key_t(ptttype.SHM_KEY), ptttype.USE_HUGETLB, false)
}

// AttachCheckSHM
//
// This is to have attach_check_shm (checking loaded)
// Should be used after LoadUHash (shmctl init) is done.
// Should be used only once in the beginning of the program.
func AttachCheckSHM() (err error) {
	err = AttachSHM()
	if err != nil {
		return err
	}

	loaded := Shm.Shm.Loaded
	if loaded == 0 {
		return ErrShmNotLoaded
	}

	// line: 135
	// commit: 6bdd36898bde207683a441cdffe2981e95de5b20
	if Shm.Shm.BTouchTime == 0 {
		Shm.Shm.BTouchTime = 1
	}

	// XXX line: 137 skip setting bcache because there is no direct-ptr for bcache for now.

	// line: 139
	if Shm.Shm.PTouchTime == 0 {
		Shm.Shm.PTouchTime = 1
	}

	// line: 142
	if Shm.Shm.FTouchTime == 0 {
		Shm.Shm.FTouchTime = 1
	}
	return nil
}
