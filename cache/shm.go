package cache

import (
	"encoding/binary"
	"reflect"
	"runtime/debug"
	"sort"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/shm"
	"github.com/Ptt-official-app/go-pttbbs/types"

	log "github.com/sirupsen/logrus"
)

type Shm struct {
	Shmid   int
	IsNew   bool
	Shmaddr unsafe.Pointer

	Raw SHMRaw // dummy variable to get the offset and size of the shm-fields.

	Shm *SHMRaw
}

func NewShm(key types.Key_t, isUseHugeTlb bool, isCreate bool) (retShm *Shm, err error) {
	shmid := int(0)
	var shmaddr unsafe.Pointer
	isNew := false

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
			return nil, err
		}
	} else {
		shmid, shmaddr, err = shm.OpenShm(key, size, isUseHugeTlb)
		if err != nil {
			return nil, err
		}
	}

	retShm = &Shm{
		Shmid:   shmid,
		IsNew:   isNew,
		Shmaddr: shmaddr,
		Shm:     (*SHMRaw)(shmaddr),
	}

	if isNew {
		retShm.Shm.Version = int32(SHM_VERSION)
		retShm.Shm.Size = int32(SHM_RAW_SZ)
		retShm.Shm.Number = int32(0)
		retShm.Shm.Loaded = int32(0)
	}

	// version and size should be fixed.
	retShm.Raw.Version = retShm.Shm.Version
	retShm.Raw.Size = retShm.Shm.Size

	// verify version
	if retShm.Raw.Version != SHM_VERSION {
		log.Errorf("cache.NewSHM: version not match: key: %v Shm.Raw.Version: %v SHM_VERSION: %v isCreate: %v isNew: %v", key, retShm.Raw.Version, SHM_VERSION, isCreate, isNew)
		debug.PrintStack()
		_ = CloseSHM()
		return nil, ErrShmVersion
	}

	if retShm.Raw.Size != int32(SHM_RAW_SZ) {
		log.Warnf("cache.NewSHM: size not match (version matched): key: %v Shm.Raw.Size: %v SHM_RAW_SZ: %v size: %v isCreate: %v isNew: %v", key, retShm.Raw.Size, SHM_RAW_SZ, size, isCreate, isNew)

		_ = CloseSHM()
		return nil, ErrShmSize
	}

	log.Infof("cache.NewSHM: shm created: key: %v shmid: %v shmaddr: %v size: %v isNew: %v", key, retShm.Shmid, retShm.Shmaddr, retShm.Raw.Size, isNew)

	return retShm, nil
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
func (s *Shm) Close() error {
	if !IsTest {
		return ErrInvalidOp
	}

	return shm.CloseShm(s.Shmid, s.Shmaddr)
}

func (s *Shm) Reset() {
	if !IsTest {
		return
	}

	const sz = SHM_RAW_SZ - uintptr(types.INT32_SZ*2)
	ptrBytes := (*[sz]byte)(unsafe.Pointer(&EMPTY_SHM_RAW.Userid))
	shmBytes := (*[SHM_RAW_SZ]byte)(unsafe.Pointer(&s.Shm.Userid))
	copy(shmBytes[:], ptrBytes[:])
}

func (s *Shm) GetBNumber() (bnumber int32) {
	return s.Shm.BNumber
}

func (s *Shm) QsortCmpBoardName() {
	bnumber := ptttype.BidInStore(s.GetBNumber())
	for i := ptttype.BidInStore(0); i < bnumber; i++ {
		s.Shm.BSorted[ptttype.BSORT_BY_NAME][i] = i
	}
	sort.Sort(shmBoardByName(s.Shm.BSorted[ptttype.BSORT_BY_NAME][:bnumber]))
}

func (s *Shm) QsortCmpBoardClass() {
	bnumber := ptttype.BidInStore(s.GetBNumber())
	for i := ptttype.BidInStore(0); i < bnumber; i++ {
		s.Shm.BSorted[ptttype.BSORT_BY_CLASS][i] = i
	}
	sort.Sort(shmBoardByClass(s.Shm.BSorted[ptttype.BSORT_BY_CLASS][:bnumber]))
}

func (s *Shm) CheckMaxUser() {
	utmpnumber := s.Shm.UTMPNumber
	maxuser := s.Shm.MaxUser
	if maxuser < utmpnumber {
		nowTS := types.NowTS()
		s.Shm.MaxUser = utmpnumber
		s.Shm.MaxTime = nowTS
	}
}
