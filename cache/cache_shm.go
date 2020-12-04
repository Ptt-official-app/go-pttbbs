package cache

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

//AttachSHM
//
//This is to have attach_shm (no checking loaded), not attach_SHM (checking loaded)
//Should be used after LoadUHash (shmctl init) is done.
//Should be used only once in the beginning of the program.
func AttachSHM() error {
	if Shm != nil {
		return nil
	}

	return NewSHM(types.Key_t(ptttype.SHM_KEY), ptttype.USE_HUGETLB, false)
}

//AttachCheckSHM
//
//This is to have attach_check_shm (checking loaded)
//Should be used after LoadUHash (shmctl init) is done.
//Should be used only once in the beginning of the program.
func AttachCheckSHM() (err error) {
	err = AttachSHM()
	if err != nil {
		return err
	}

	loaded := int32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Loaded),
		types.INT32_SZ,
		unsafe.Pointer(&loaded),
	)
	if loaded == 0 {
		return ErrShmNotLoaded
	}

	// line: 135
	// commit: 6bdd36898bde207683a441cdffe2981e95de5b20
	btouchTime := types.Time4(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.BTouchTime),
		types.TIME4_SZ,
		unsafe.Pointer(&btouchTime),
	)
	if btouchTime == 0 {
		btouchTime = types.Time4(1)
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.BTouchTime),
			types.TIME4_SZ,
			unsafe.Pointer(&btouchTime),
		)
	}

	//XXX line: 137 skip setting bcache because there is no direct-ptr for bcache for now.

	// line: 139
	ptouchTime := types.Time4(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.PTouchTime),
		types.TIME4_SZ,
		unsafe.Pointer(&ptouchTime),
	)
	if ptouchTime == 0 {
		ptouchTime = 1
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.PTouchTime),
			types.TIME4_SZ,
			unsafe.Pointer(&ptouchTime),
		)
	}

	// line: 142
	ftouchTime := types.Time4(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.PTouchTime),
		types.TIME4_SZ,
		unsafe.Pointer(&ftouchTime),
	)
	if ftouchTime == 0 {
		ftouchTime = 1
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.FTouchTime),
			types.TIME4_SZ,
			unsafe.Pointer(&ftouchTime),
		)
	}
	return nil
}
