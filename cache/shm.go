package cache

import (
    "encoding/binary"
    "reflect"
    "runtime/debug"
    "unsafe"

    "github.com/Ptt-official-app/go-pttbbs/ptttype"
    "github.com/Ptt-official-app/go-pttbbs/shm"
    "github.com/Ptt-official-app/go-pttbbs/types"

    log "github.com/sirupsen/logrus"
)

type SHM struct {
    Shmid   int
    IsNew   bool
    Shmaddr unsafe.Pointer

    Raw SHMRaw //dummy variable to get the offset and size of the shm-fields.
}

//NewSHM
//
//This is to init SHM with Version and Size checked.
func NewSHM(key types.Key_t, isUseHugeTlb bool, isCreate bool) error {
    if Shm != nil {
        return ErrShmAlreadyInit
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
    debugShm()

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
    }

    if isNew {
        in_version := SHM_VERSION
        in_size := int32(SHM_RAW_SZ)
        in_number := int32(0)
        in_loaded := int32(0)
        Shm.WriteAt(
            unsafe.Offsetof(Shm.Raw.Version),
            unsafe.Sizeof(Shm.Raw.Version),
            unsafe.Pointer(&in_version),
        )
        Shm.WriteAt(
            unsafe.Offsetof(Shm.Raw.Size),
            unsafe.Sizeof(Shm.Raw.Size),
            unsafe.Pointer(&in_size),
        )
        Shm.WriteAt(
            unsafe.Offsetof(Shm.Raw.Number),
            unsafe.Sizeof(Shm.Raw.Number),
            unsafe.Pointer(&in_number),
        )
        Shm.WriteAt(
            unsafe.Offsetof(Shm.Raw.Loaded),
            unsafe.Sizeof(Shm.Raw.Loaded),
            unsafe.Pointer(&in_loaded),
        )
    }

    // version and size should be fixed.
    Shm.ReadAt(
        unsafe.Offsetof(Shm.Raw.Version),
        unsafe.Sizeof(Shm.Raw.Version),
        unsafe.Pointer(&Shm.Raw.Version),
    )
    Shm.ReadAt(
        unsafe.Offsetof(Shm.Raw.Size),
        unsafe.Sizeof(Shm.Raw.Size),
        unsafe.Pointer(&Shm.Raw.Size),
    )

    Shm.SetBCACHEPTR(
        unsafe.Offsetof(Shm.Raw.BCache),
    )

    // verify version
    if Shm.Raw.Version != SHM_VERSION {
        log.Errorf("cache.NewSHM: version not match: key: %v Shm.Raw.Version: %v SHM_VERSION: %v isCreate: %v isNew: %v", key, Shm.Raw.Version, SHM_VERSION, isCreate, isNew)
        debug.PrintStack()
        CloseSHM()
        return ErrShmVersion
    }
    if Shm.Raw.Size != int32(SHM_RAW_SZ) {
        log.Warnf("cache.NewSHM: size not match (version matched): key: %v Shm.Raw.Size: %v SHM_RAW_SZ: %v size: %v isCreate: %v isNew: %v", key, Shm.Raw.Size, SHM_RAW_SZ, size, isCreate, isNew)

        CloseSHM()
        return ErrShmSize
    }

    if isCreate && !isNew {
        log.Warnf("cache.NewSHM: is expected to create, but not: key: %v", key)
    }

    log.Infof("cache.NewSHM: shm created: key: %v size: %v isNew: %v", key, Shm.Raw.Size, isNew)

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

//Close
//
//XXX [WARNING] know what you are doing before using Close!.
//This is to be able to close the shared mem for the completeness of the mem-usage.
//However, in production, we create shm without the need of closing the shm.
//
//We simply use ipcrm to delete the shm if necessary.
//
//Currently used only in test.
func CloseSHM() error {
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

    log.Infof("cache.CloseSHM: done")

    return nil
}

//Close
//
//XXX [WARNING] know what you are doing before using Close!.
//This is to be able to close the shared mem for the completeness of the mem-usage.
//However, in production, we create shm without the need of closing the shm.
//
//We simply use ipcrm to delete the shm if necessary.
//
//Currently used only in test.
func (s *SHM) Close() error {
    if !IsTest {
        return ErrInvalidOp
    }
    return shm.CloseShm(s.Shmid, s.Shmaddr)
}

//ReadAt
//
//Require precalculated offset and size and outptr to efficiently get the data.
//See tests for exact usage.
//[!!!] If we are reading from the array, make sure that have unit-size * n in the size.
//
//Params
//  offsetOfSHMRawComponent: offset from SHMRaw
//  size: size of the variable, usually can be referred from SHMRaw
//        [!!!]If we are reading from the array, make sure that have unit-size * n in the size.
//  outptr: the ptr of the object to read.
func (s *SHM) ReadAt(offsetOfSHMRawComponent uintptr, size uintptr, outptr unsafe.Pointer) {
    shm.ReadAt(s.Shmaddr, int(offsetOfSHMRawComponent), size, outptr)
}

//WriteAt
//
//Require recalculated offset and size and outptr to efficiently get the data.
//See tests for exact usage.
//[!!!]If we are reading from the array, make sure that have unit-size * n in the size.
//
//Params
//  offsetOfSHMRawComponent: offset from SHMRaw
//  size: size of the variable
//        [!!!]If we are reading from the array, make sure that have unit-size * n in the size.
//  inptr: the ptr of the object to write.
func (s *SHM) WriteAt(offsetOfSHMRawComponent uintptr, size uintptr, inptr unsafe.Pointer) {
    shm.WriteAt(s.Shmaddr, int(offsetOfSHMRawComponent), size, inptr)
}

func (s *SHM) SetOrUint32(offsetOfSHMRawComponent uintptr, flag uint32) {
    shm.SetOrUint32(s.Shmaddr, int(offsetOfSHMRawComponent), flag)
}

func (s *SHM) IncUint32(offsetOfSHMRawComponent uintptr) {
    shm.IncUint32(s.Shmaddr, int(offsetOfSHMRawComponent))
}

func (s *SHM) Memset(offsetOfSHMRawComponent uintptr, c byte, size uintptr) {
    shm.Memset(s.Shmaddr, int(offsetOfSHMRawComponent), c, size)
}

func (s *SHM) InnerSetInt32(offsetSrc uintptr, offsetDst uintptr) {
    shm.InnerSetInt32(s.Shmaddr, int(offsetSrc), int(offsetDst))
}

func (s *SHM) Memcmp(offsetOfSHMRawComponent uintptr, size uintptr, cmpptr unsafe.Pointer) int {
    return shm.Memcmp(s.Shmaddr, int(offsetOfSHMRawComponent), size, cmpptr)
}

//SetBCACHEPTR
//
//!!!Required in NewSHM (and should be set only once in NewSHM)
func (s *SHM) SetBCACHEPTR(offsetOfSHMRawComponent uintptr) {
    shm.SetBCACHEPTR(s.Shmaddr, int(offsetOfSHMRawComponent))
}

func (s *SHM) GetBNumber() (bnumber int32) {
    s.ReadAt(
        unsafe.Offsetof(s.Raw.BNumber),
        types.TIME4_SZ,
        unsafe.Pointer(&bnumber),
    )
    return
}

func (s *SHM) QsortCmpBoardName() {
    bnumber := s.GetBNumber()

    const bsorted0sz = unsafe.Sizeof(s.Raw.BSorted[0])
    offsetBsorted := unsafe.Offsetof(s.Raw.BSorted) + bsorted0sz*uintptr(ptttype.BSORT_BY_NAME)
    shm.QsortCmpBoardName(s.Shmaddr, int(offsetBsorted), uint32(bnumber))
}

func (s *SHM) QsortCmpBoardClass() {
    bnumber := s.GetBNumber()

    const bsorted0sz = unsafe.Sizeof(s.Raw.BSorted[0])
    offsetBsorted := unsafe.Offsetof(s.Raw.BSorted) + bsorted0sz*uintptr(ptttype.BSORT_BY_CLASS)
    shm.QsortCmpBoardClass(s.Shmaddr, int(offsetBsorted), uint32(bnumber))
}
