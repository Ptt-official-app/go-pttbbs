package cache

import (
	"encoding/binary"
	"reflect"
	"runtime/debug"
	"sort"
	"time"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/shm"
	"github.com/Ptt-official-app/go-pttbbs/types"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type SHM struct {
	Shmid   int
	IsNew   bool
	Shmaddr unsafe.Pointer

	Raw SHMRaw // dummy variable to get the offset and size of the shm-fields.

	Shm *SHMRaw
}

// NewSHM
//
// This is to init SHM with Version and Size checked.
func NewSHM(key types.Key_t, isUseHugeTlb bool, isCreate bool) error {
	if Shm != nil {
		return nil
	}

	shmid := int(0)
	var shmaddr unsafe.Pointer
	isNew := false
	var err error

	// pttstruct.h line: 616
	SHMSIZE := types.Size_t(SHM_RAW_SZ)
	if ptttype.SHMALIGNEDSIZE != 0 {
		SHMSIZE = types.Size_t((int(SHM_RAW_SZ)/(ptttype.SHMALIGNEDSIZE) + 1) * ptttype.SHMALIGNEDSIZE)
	}

	log.Infof("cache.NewSHM: SHMSIZE: %v SHM_RAW_SZ: %v SHMALIGNEDSIZE: %v", SHMSIZE, SHM_RAW_SZ, ptttype.SHMALIGNEDSIZE)
	if !IsTest {
		debugShm()
	}

	size := SHMSIZE

	if isCreate {
		shmid, shmaddr, isNew, err = shm.CreateShm(key, size, isUseHugeTlb)
		if err != nil {
			return err
		}
	} else {
		shmid, shmaddr, err = shm.OpenShm(key, size, isUseHugeTlb)
		if err != nil {
			return err
		}
	}

	Shm = &SHM{
		Shmid:   shmid,
		IsNew:   isNew,
		Shmaddr: shmaddr,
		Shm:     (*SHMRaw)(shmaddr),
	}

	if isNew {
		in_version := SHM_VERSION
		in_size := int32(SHM_RAW_SZ)
		in_number := int32(0)
		in_loaded := int32(0)
		Shm.Shm.Version = int32(in_version)
		Shm.Shm.Size = in_size
		Shm.Shm.Number = in_number
		Shm.Shm.Loaded = in_loaded
	}

	// version and size should be fixed.
	Shm.Raw.Version = Shm.Shm.Version
	Shm.Raw.Size = Shm.Shm.Size

	// verify version
	logrus.Infof("cache.NewSHM: shmid: %v isNew: %v shmaddr: %v Raw.Version: %v SHM_VERSION: %v", shmid, isNew, shmaddr, Shm.Raw.Version, SHM_VERSION)
	if Shm.Raw.Version != SHM_VERSION {
		log.Errorf("cache.NewSHM: version not match: key: %v Shm.Raw.Version: %v SHM_VERSION: %v isCreate: %v isNew: %v", key, Shm.Raw.Version, SHM_VERSION, isCreate, isNew)
		debug.PrintStack()
		_ = CloseSHM()
		return ErrShmVersion
	}
	logrus.Infof("cache.NewSHM: Raw.Size: %v SHM_RAW_SZ: %v", Shm.Raw.Size, SHM_RAW_SZ)
	if Shm.Raw.Size != int32(SHM_RAW_SZ) {
		log.Warnf("cache.NewSHM: size not match (version matched): key: %v Shm.Raw.Size: %v SHM_RAW_SZ: %v size: %v isCreate: %v isNew: %v", key, Shm.Raw.Size, SHM_RAW_SZ, size, isCreate, isNew)

		_ = CloseSHM()
		return ErrShmSize
	}

	if isCreate && !isNew {
		log.Warnf("cache.NewSHM: is expected to create, but not: key: %v", key)
	}

	log.Infof("cache.NewSHM: shm created: key: %v shmid: %v shmaddr: %v size: %v isNew: %v", key, Shm.Shmid, Shm.Shmaddr, Shm.Raw.Size, isNew)

	return nil
}

func debugShm() {
	val := reflect.ValueOf(EMPTY_SHM_RAW)
	numField := val.NumField()
	for i := 0; i < numField; i++ {
		field := val.Type().Field(i)
		fieldName := field.Name
		offset := field.Offset

		value := val.Field(i).Interface()
		theSize := binary.Size(value)

		log.Debugf("cache.Shm.Raw: (%v/%v) %v: offset: %v size: %v", i, numField, fieldName, offset, theSize)
	}
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

	if Shm == nil {
		// Already Closed
		log.Errorf("cache.CloseSHM: already closed")
		return ErrShmNotInit
	}

	err := Shm.Close()
	if err != nil {
		log.Errorf("cache.CloseSHM: unable to close: e: %v", err)
		return err
	}

	Shm = nil

	time.Sleep(3 * time.Millisecond)

	log.Infof("cache.CloseSHM: done")

	return nil
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
func (s *SHM) Close() error {
	if !IsTest {
		return ErrInvalidOp
	}

	return shm.CloseShm(s.Shmid, s.Shmaddr)
}

func (s *SHM) Reset() {
	if !IsTest {
		return
	}

	const sz = SHM_RAW_SZ - uintptr(types.INT32_SZ*2)
	ptrBytes := (*[sz]byte)(unsafe.Pointer(&EMPTY_SHM_RAW.Userid))
	shmBytes := (*[SHM_RAW_SZ]byte)(unsafe.Pointer(&s.Shm.Userid))
	copy(shmBytes[:], ptrBytes[:])
}

func (s *SHM) GetBNumber() (bnumber int32) {
	return s.Shm.BNumber
}

func (s *SHM) QsortCmpBoardName() {
	bnumber := ptttype.BidInStore(s.GetBNumber())
	for i := ptttype.BidInStore(0); i < bnumber; i++ {
		s.Shm.BSorted[ptttype.BSORT_BY_NAME][i] = i
	}
	sort.Sort(shmBoardByName(s.Shm.BSorted[ptttype.BSORT_BY_NAME][:bnumber]))
}

func (s *SHM) QsortCmpBoardClass() {
	bnumber := ptttype.BidInStore(s.GetBNumber())
	for i := ptttype.BidInStore(0); i < bnumber; i++ {
		s.Shm.BSorted[ptttype.BSORT_BY_CLASS][i] = i
	}
	sort.Sort(shmBoardByClass(s.Shm.BSorted[ptttype.BSORT_BY_CLASS][:bnumber]))
}

func (s *SHM) CheckMaxUser() {
	utmpnumber := s.Shm.UTMPNumber
	maxuser := s.Shm.MaxUser
	if maxuser < utmpnumber {
		nowTS := types.NowTS()
		s.Shm.MaxUser = utmpnumber
		s.Shm.MaxTime = nowTS
	}
}
