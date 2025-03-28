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
	if SHM != nil {
		return nil
	}

	if types.IS_ALL_GUEST {
		return nil
	}

	return Init(types.Key_t(ptttype.SHM_KEY), ptttype.USE_HUGETLB, false)
}

// AttachCheckSHM
//
// This is to have attach_check_shm (checking loaded)
// Should be used after LoadUHash (shmctl init) is done.
// Should be used only once in the beginning of the program.
func AttachCheckSHM() (err error) {
	if types.IS_ALL_GUEST {
		return nil
	}

	err = AttachSHM()
	if err != nil {
		return err
	}

	loaded := SHM.Shm.Loaded
	if loaded == 0 {
		return ErrShmNotLoaded
	}

	// line: 135
	// commit: 6bdd36898bde207683a441cdffe2981e95de5b20
	if SHM.Shm.BTouchTime == 0 {
		SHM.Shm.BTouchTime = 1
	}

	// XXX line: 137 skip setting bcache because there is no direct-ptr for bcache for now.

	// line: 139
	if SHM.Shm.PTouchTime == 0 {
		SHM.Shm.PTouchTime = 1
	}

	// line: 142
	if SHM.Shm.FTouchTime == 0 {
		SHM.Shm.FTouchTime = 1
	}
	return nil
}
