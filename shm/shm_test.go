package shm

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

type testStruct struct {
	A int32
	B testStruct2
}

const TEST_STRUCT_SZ = unsafe.Sizeof(testStruct{})

type testStruct2 struct {
	C [10]uint8
}

func TestCreateShm(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		key           types.Key_t
		size          types.Size_t
		is_usehugetlb bool
	}
	tests := []struct {
		name        string
		args        args
		wantShmid   int
		wantShmaddr unsafe.Pointer
		wantIsNew   bool
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name:      "first-good-shm-id",
			args:      args{testShmKey, 100, false},
			wantShmid: 0,
			wantIsNew: true,
		},
		{
			name:      "not first-good-shm-id",
			args:      args{testShmKey, 100, false},
			wantShmid: 0,
			wantIsNew: false,
		},
	}

	firstGoodShmID := 0
	var firstGoodShmaddr unsafe.Pointer
	defer CloseShm(firstGoodShmID, firstGoodShmaddr)
	for idx, tt := range tests {
		gotShmid, gotShmaddr, gotIsNew, err := CreateShm(tt.args.key, tt.args.size, tt.args.is_usehugetlb)
		log.Infof("(%v/%v): after CreateShm: gotShmid: %v gotShmaddr: %v gotIsNew: %v e: %v", idx, tt.name, gotShmid, gotShmaddr, gotIsNew, err)

		if (err != nil) != tt.wantErr {
			t.Errorf("(%v/%v): CreateShm() error = %v, wantErr %v", idx, tt.name, err, tt.wantErr)
			return
		}

		log.Infof("shm_test.CreateShm: to check firstGoodShmID: %v %v", firstGoodShmID, firstGoodShmaddr)
		if firstGoodShmID == 0 {
			firstGoodShmID = gotShmid
			firstGoodShmaddr = gotShmaddr
		}

		if tt.wantShmid != 0 && gotShmid != tt.wantShmid {
			t.Errorf("%v: CreateShm() gotShmid = %v, expected %v", tt.name, gotShmid, tt.wantShmid)
		}
		if gotIsNew != tt.wantIsNew {
			t.Errorf("%v: CreateShm() gotIsNew = %v, expected %v", tt.name, gotIsNew, tt.wantIsNew)
		}
	}
}

func TestCloseShm(t *testing.T) {
	setupTest()
	defer teardownTest()

	gotShmid, gotShmaddr, _, _ := CreateShm(testShmKey, 100, false)
	log.Infof("TestCloseShm: after CreateShm: gotShmid: %v gotShmaddr: %v", gotShmid, gotShmaddr)

	type args struct {
		shmid   int
		shmaddr unsafe.Pointer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{gotShmid, gotShmaddr},
		},
		{
			args:    args{gotShmid, gotShmaddr},
			wantErr: true,
		},
	}
	for idx, tt := range tests {
		err := CloseShm(tt.args.shmid, tt.args.shmaddr)
		if (err != nil) != tt.wantErr {
			t.Errorf("(%v/%v) CloseShm: e: %v wantErr: %v", idx, tt.name, err, tt.wantErr)
		}
	}
}

func TestOpenShm(t *testing.T) {
	setupTest()
	defer teardownTest()

	shmid, shmaddr, _, _ := CreateShm(testShmKey, 100, false)
	defer CloseShm(shmid, shmaddr)

	type args struct {
		key           types.Key_t
		size          types.Size_t
		is_usehugetlb bool
	}
	tests := []struct {
		name        string
		args        args
		wantShmid   int
		wantShmaddr unsafe.Pointer
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name:      "same size",
			args:      args{testShmKey, 100, false},
			wantShmid: shmid,
		},
		{
			name:      "same size 2",
			args:      args{testShmKey, 100, false},
			wantShmid: shmid,
		},
		{
			name:      "diff size",
			args:      args{testShmKey, 500, false},
			wantShmid: -1,
			wantErr:   true,
		},
		{
			name:      "diff key",
			args:      args{testShmKey + 1, 500, false},
			wantShmid: -1,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotShmid, _, err := OpenShm(tt.args.key, tt.args.size, tt.args.is_usehugetlb)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenShm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotShmid != tt.wantShmid {
				t.Errorf("OpenShm() gotShmid = %v, expected %v", gotShmid, tt.wantShmid)
			}
		})
	}

}

func TestReadAt(t *testing.T) {
	setupTest()
	defer teardownTest()

	shmid, shmaddr, _, _ := CreateShm(testShmKey, 40, false)
	defer CloseShm(shmid, shmaddr)

	_, shmaddr2, _ := OpenShm(testShmKey, 40, false)

	empty := [40]byte{}
	emptyptr := unsafe.Pointer(&empty)
	WriteAt(shmaddr, 0, 40, emptyptr)

	test1 := &testStruct{}
	test1.A = 10
	copy(test1.B.C[:], []byte("0123456789"))
	test1ptr := unsafe.Pointer(test1)
	WriteAt(shmaddr, 12, TEST_STRUCT_SZ, test1ptr)
	log.Infof("TEST_STRUCT_SZ: %v", TEST_STRUCT_SZ)

	test2 := &testStruct{}
	test2ptr := unsafe.Pointer(test2)

	test3 := &[40]byte{}
	test3ptr := unsafe.Pointer(test3)

	test4 := &[40]byte{}
	test4ptr := unsafe.Pointer(test4)

	want3 := &[40]byte{}
	want3[12] = 10
	copy(want3[12+4:], []byte("0123456789"))

	type args struct {
		shmaddr unsafe.Pointer
		offset  int
		size    uintptr
		outptr  unsafe.Pointer
	}
	tests := []struct {
		name     string
		args     args
		read     interface{}
		expected interface{}
	}{
		// TODO: Add test cases.
		{
			name:     "read test2 (testStruct)",
			args:     args{shmaddr: shmaddr, offset: 12, size: TEST_STRUCT_SZ, outptr: test2ptr},
			read:     test2,
			expected: test1,
		},
		{
			name:     "read test3 (bytes)",
			args:     args{shmaddr: shmaddr, offset: 0, size: 40, outptr: test3ptr},
			read:     test3,
			expected: want3,
		},
		{
			name:     "read test4 (bytes) from shmaddr2",
			args:     args{shmaddr: shmaddr2, offset: 0, size: 40, outptr: test4ptr},
			read:     test4,
			expected: want3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReadAt(tt.args.shmaddr, tt.args.offset, tt.args.size, tt.args.outptr)

			if !reflect.DeepEqual(tt.read, tt.expected) {
				t.Errorf("ReadAt() read: %v expected: %v", tt.read, tt.expected)
			}
		})
	}

}

func TestWriteAt(t *testing.T) {
	setupTest()
	defer teardownTest()

	shmid, shmaddr, _, _ := CreateShm(testShmKey, 40, false)
	defer CloseShm(shmid, shmaddr)

	empty := [40]byte{}
	emptyptr := unsafe.Pointer(&empty)

	empty1 := [40]byte{}

	WriteAt(shmaddr, 0, 40, emptyptr)

	read := [40]byte{}
	readptr := unsafe.Pointer(&read)

	test1 := &testStruct{}
	test1.A = 10
	copy(test1.B.C[:], []byte("0123456789"))
	test1ptr := unsafe.Pointer(test1)

	want1 := [40]byte{}
	want1[0] = 10
	copy(want1[0+4:], []byte("0123456789"))

	want2 := [40]byte{}
	copy(want2[:12], want1[:])
	want2[12] = 10
	copy(want2[12+4:], []byte("0123456789"))

	type args struct {
		shmaddr unsafe.Pointer
		offset  int
		size    uintptr
		inptr   unsafe.Pointer
	}
	tests := []struct {
		name     string
		args     args
		expected [40]byte
	}{
		// TODO: Add test cases.
		{
			name:     "empty (1)",
			args:     args{shmaddr, 0, 40, emptyptr},
			expected: empty1,
		},
		{
			name:     "test",
			args:     args{shmaddr, 0, TEST_STRUCT_SZ, test1ptr},
			expected: want1,
		},
		{
			name:     "test-want2",
			args:     args{shmaddr, 12, TEST_STRUCT_SZ, test1ptr},
			expected: want2,
		},
		{
			name:     "empty (2)",
			args:     args{shmaddr, 0, 40, emptyptr},
			expected: empty1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteAt(tt.args.shmaddr, tt.args.offset, tt.args.size, tt.args.inptr)

			ReadAt(shmaddr, 0, 40, readptr)
			if !reflect.DeepEqual(read, tt.expected) {
				t.Errorf("WriteAt() = %v expected: %v", read, tt.expected)
			}

		})
	}

}

func TestIncUint32(t *testing.T) {
	setupTest()
	defer teardownTest()

	shmid, shmaddr, _, _ := CreateShm(testShmKey, 40, false)
	defer CloseShm(shmid, shmaddr)

	empty := [40]byte{}
	emptyptr := unsafe.Pointer(&empty)
	WriteAt(shmaddr, 0, 40, emptyptr)

	type args struct {
		shmaddr unsafe.Pointer
		offset  int
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
	}{
		// TODO: Add test cases.
		{
			args:     args{shmaddr: shmaddr, offset: 1},
			expected: 1,
		},
		{
			args:     args{shmaddr: shmaddr, offset: 1},
			expected: 2,
		},
		{
			args:     args{shmaddr: shmaddr, offset: 2},
			expected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IncUint32(tt.args.shmaddr, tt.args.offset)

			out := uint32(0)
			ReadAt(tt.args.shmaddr, tt.args.offset, unsafe.Sizeof(out), unsafe.Pointer(&out))
			if !reflect.DeepEqual(out, tt.expected) {
				t.Errorf("IncUint32() out: %v expected: %v", out, tt.expected)
			}
		})
	}
}

func TestSetOrUint32(t *testing.T) {
	setupTest()
	defer teardownTest()

	shmid, shmaddr, _, _ := CreateShm(testShmKey, 40, false)
	defer CloseShm(shmid, shmaddr)

	type args struct {
		shmaddr unsafe.Pointer
		offset  int
		flag    uint32
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
	}{
		// TODO: Add test cases.
		{
			args:     args{shmaddr: shmaddr, offset: 3, flag: 0x00000005},
			expected: 0x00000005,
		},
		{
			args:     args{shmaddr: shmaddr, offset: 3, flag: 0x00000050},
			expected: 0x00000055,
		},
		{
			args:     args{shmaddr: shmaddr, offset: 3, flag: 0x00000014},
			expected: 0x00000055,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetOrUint32(tt.args.shmaddr, tt.args.offset, tt.args.flag)

			out := uint32(0)
			ReadAt(
				tt.args.shmaddr,
				tt.args.offset,
				unsafe.Sizeof(out),
				unsafe.Pointer(&out),
			)
			if out != tt.expected {
				t.Errorf("SetOrUint32() out: %v expected: %v", out, tt.expected)
			}
		})
	}
}

func TestInnerSetInt32(t *testing.T) {
	setupTest()
	defer teardownTest()

	shmid, shmaddr, _, _ := CreateShm(testShmKey, 40, false)
	defer CloseShm(shmid, shmaddr)

	in := int32(8)
	WriteAt(
		shmaddr,
		12,
		types.INT32_SZ,
		unsafe.Pointer(&in),
	)

	type args struct {
		shmaddr   unsafe.Pointer
		offsetSrc int
		offsetDst int
	}
	tests := []struct {
		name     string
		args     args
		expected int32
	}{
		// TODO: Add test cases.
		{
			args:     args{shmaddr: shmaddr, offsetSrc: 12, offsetDst: 10},
			expected: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InnerSetInt32(tt.args.shmaddr, tt.args.offsetSrc, tt.args.offsetDst)

			out := int32(0)
			ReadAt(
				tt.args.shmaddr,
				tt.args.offsetDst,
				unsafe.Sizeof(out),
				unsafe.Pointer(&out),
			)
			if out != tt.expected {
				t.Errorf("InnerSetInt32(): %v expected: %v", out, tt.expected)
			}
		})
	}
}

func TestMemset(t *testing.T) {
	setupTest()
	defer teardownTest()

	shmid, shmaddr, _, _ := CreateShm(testShmKey, 40, false)
	defer CloseShm(shmid, shmaddr)

	type args struct {
		shmaddr unsafe.Pointer
		offset  int
		c       byte
		size    uintptr
	}
	tests := []struct {
		name     string
		args     args
		expected []byte
	}{
		// TODO: Add test cases.
		{
			args: args{shmaddr: shmaddr, offset: 10, c: 'a', size: 5},
			expected: []byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				'a', 'a', 'a', 'a', 'a', 0, 0, 0, 0, 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Memset(tt.args.shmaddr, tt.args.offset, tt.args.c, tt.args.size)

			out := [20]byte{}
			ReadAt(
				shmaddr,
				0,
				20,
				unsafe.Pointer(&out),
			)

			if !reflect.DeepEqual(out[:], tt.expected) {
				t.Errorf("Memset: %v expected: %v", out, tt.expected)
			}
		})
	}
}

func TestQsortCmpBoardName(t *testing.T) {
	setupTest()
	defer teardownTest()

	shmid, shmaddr, _, _ := CreateShm(testShmKey, 5000, false)
	defer CloseShm(shmid, shmaddr)

	SetBCACHEPTR(shmaddr, 1200)

	brdname := ptttype.BoardID_t{}
	brdtitle := ptttype.BoardTitle_t{}
	log.Infof("QsortCmpBoardName: BOARD_HEADER_RAW_SZ: %v brdname_sz: %v brd_title_sz: %v shmaddr: %v", ptttype.BOARD_HEADER_RAW_SZ, unsafe.Sizeof(ptttype.BoardID_t{}), unsafe.Sizeof(ptttype.BoardTitle_t{}), shmaddr)
	copy(brdname[:], "test_board3")
	copy(brdtitle[:], "心情 title3")

	outBoardname := ptttype.BoardID_t{}
	WriteAt(
		shmaddr,
		1200,
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&brdname),
	)
	ReadAt(
		shmaddr,
		1200,
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&outBoardname),
	)
	if !reflect.DeepEqual(brdname, outBoardname) {
		t.Errorf("unable to write: brdname: %v outBoardname: %v\n", brdname, outBoardname)
	}
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
		ptttype.BOARD_TITLE_SZ,
		unsafe.Pointer(&brdtitle),
	)

	copy(brdname[:], "test_board4")
	copy(brdtitle[:], "心情 title4")
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ),
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&brdname),
	)

	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
		ptttype.BOARD_TITLE_SZ,
		unsafe.Pointer(&brdtitle),
	)

	copy(brdname[:], "test_board1")
	copy(brdtitle[:], "站務 title2")
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*2,
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&brdname),
	)
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*2+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
		ptttype.BOARD_TITLE_SZ,
		unsafe.Pointer(&brdtitle),
	)

	copy(brdname[:], "test_board5")
	copy(brdtitle[:], "站務 title1")
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*3,
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&brdname),
	)
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*3+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
		ptttype.BOARD_TITLE_SZ,
		unsafe.Pointer(&brdtitle),
	)

	copy(brdname[:], "test_board2")
	copy(brdtitle[:], "球類 title3")
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*4,
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&brdname),
	)
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*4+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
		ptttype.BOARD_TITLE_SZ,
		unsafe.Pointer(&brdtitle),
	)

	type out_t struct {
		Name  []byte
		Title []byte
	}
	expected := []out_t{
		{
			Name:  []byte("test_board1"),
			Title: []byte("站務 title2"),
		},
		{
			Name:  []byte("test_board2"),
			Title: []byte("球類 title3"),
		},
		{
			Name:  []byte("test_board3"),
			Title: []byte("心情 title3"),
		},
		{
			Name:  []byte("test_board4"),
			Title: []byte("心情 title4"),
		},
		{
			Name:  []byte("test_board5"),
			Title: []byte("站務 title1"),
		},
	}

	type args struct {
		shmaddr unsafe.Pointer
		offset  int
		n       uint32
	}
	tests := []struct {
		name     string
		args     args
		expected []out_t
	}{
		// TODO: Add test cases.
		{
			args:     args{shmaddr: shmaddr, offset: 200, n: 5},
			expected: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			QsortCmpBoardName(tt.args.shmaddr, tt.args.offset, tt.args.n)
		})

		outName := ptttype.BoardID_t{}
		outTitle := ptttype.BoardTitle_t{}
		the_idx := int32(0)
		for i := 0; i < 5; i++ {
			ReadAt(
				tt.args.shmaddr,
				tt.args.offset+i*int(types.INT32_SZ),
				types.INT32_SZ,
				unsafe.Pointer(&the_idx),
			)
			ReadAt(
				tt.args.shmaddr,
				1200+int(ptttype.BOARD_HEADER_RAW_SZ)*int(the_idx)+int(ptttype.BOARD_HEADER_BRDNAME_OFFSET),
				ptttype.BOARD_ID_SZ,
				unsafe.Pointer(&outName),
			)

			theBytes := types.CstrToBytes(outName[:])
			if !reflect.DeepEqual(theBytes, tt.expected[i].Name[:]) {
				t.Errorf("QsortCmpBoardName (name): (%v/5) (%v) expected: (%v)", i, string(theBytes), string(tt.expected[i].Name[:]))
			}

			ReadAt(
				tt.args.shmaddr,
				1200+int(ptttype.BOARD_HEADER_RAW_SZ)*int(the_idx)+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
				ptttype.BOARD_TITLE_SZ,
				unsafe.Pointer(&outTitle),
			)

			theBytes = types.CstrToBytes(outTitle[:])
			if !reflect.DeepEqual(theBytes, tt.expected[i].Title[:]) {
				t.Errorf("QsortCmpBoardName: (title): (%v/5) %v expected: %v", i, string(theBytes), string(tt.expected[i].Title[:]))
			}
		}
	}
}

func TestQsortCmpBoardClass(t *testing.T) {
	setupTest()
	defer teardownTest()

	shmid, shmaddr, _, _ := CreateShm(testShmKey, 5000, false)
	defer CloseShm(shmid, shmaddr)

	SetBCACHEPTR(shmaddr, 1200)

	brdname := ptttype.BoardID_t{}
	brdtitle := ptttype.BoardTitle_t{}
	log.Infof("QsortCmpBoardName: BOARD_HEADER_RAW_SZ: %v brdname_sz: %v brd_title_sz: %v shmaddr: %v", ptttype.BOARD_HEADER_RAW_SZ, unsafe.Sizeof(ptttype.BoardID_t{}), unsafe.Sizeof(ptttype.BoardTitle_t{}), shmaddr)
	copy(brdname[:], "test_board3")
	copy(brdtitle[:], "心情 title3")

	outBoardname := ptttype.BoardID_t{}
	WriteAt(
		shmaddr,
		1200,
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&brdname),
	)
	ReadAt(
		shmaddr,
		1200,
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&outBoardname),
	)
	if !reflect.DeepEqual(brdname, outBoardname) {
		t.Errorf("unable to write: brdname: %v outBoardname: %v\n", brdname, outBoardname)
	}
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
		ptttype.BOARD_TITLE_SZ,
		unsafe.Pointer(&brdtitle),
	)

	copy(brdname[:], "test_board4")
	copy(brdtitle[:], "心情 title4")
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ),
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&brdname),
	)

	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
		ptttype.BOARD_TITLE_SZ,
		unsafe.Pointer(&brdtitle),
	)

	copy(brdname[:], "test_board1")
	copy(brdtitle[:], "站務 title2")
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*2,
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&brdname),
	)
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*2+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
		ptttype.BOARD_TITLE_SZ,
		unsafe.Pointer(&brdtitle),
	)

	copy(brdname[:], "test_board5")
	copy(brdtitle[:], "站務 title1")
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*3,
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&brdname),
	)
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*3+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
		ptttype.BOARD_TITLE_SZ,
		unsafe.Pointer(&brdtitle),
	)

	copy(brdname[:], "test_board2")
	copy(brdtitle[:], "球類 title3")
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*4,
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(&brdname),
	)
	WriteAt(
		shmaddr,
		1200+int(ptttype.BOARD_HEADER_RAW_SZ)*4+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
		ptttype.BOARD_TITLE_SZ,
		unsafe.Pointer(&brdtitle),
	)

	type out_t struct {
		Name  []byte
		Title []byte
	}
	expected := []out_t{
		{
			Name:  []byte("test_board3"),
			Title: []byte("心情 title3"),
		},
		{
			Name:  []byte("test_board4"),
			Title: []byte("心情 title4"),
		},
		{
			Name:  []byte("test_board2"),
			Title: []byte("球類 title3"),
		},
		{
			Name:  []byte("test_board1"),
			Title: []byte("站務 title2"),
		},
		{
			Name:  []byte("test_board5"),
			Title: []byte("站務 title1"),
		},
	}

	type args struct {
		shmaddr unsafe.Pointer
		offset  int
		n       uint32
	}
	tests := []struct {
		name     string
		args     args
		expected []out_t
	}{
		// TODO: Add test cases.
		{
			args:     args{shmaddr, 200, 5},
			expected: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			QsortCmpBoardClass(tt.args.shmaddr, tt.args.offset, tt.args.n)
		})

		outName := ptttype.BoardID_t{}
		outTitle := ptttype.BoardTitle_t{}
		the_idx := int32(0)
		for i := 0; i < 5; i++ {
			ReadAt(
				tt.args.shmaddr,
				tt.args.offset+i*int(types.INT32_SZ),
				types.INT32_SZ,
				unsafe.Pointer(&the_idx),
			)
			ReadAt(
				tt.args.shmaddr,
				1200+int(ptttype.BOARD_HEADER_RAW_SZ)*int(the_idx)+int(ptttype.BOARD_HEADER_BRDNAME_OFFSET),
				ptttype.BOARD_ID_SZ,
				unsafe.Pointer(&outName),
			)

			theBytes := types.CstrToBytes(outName[:])
			if !reflect.DeepEqual(theBytes, tt.expected[i].Name[:]) {
				t.Errorf("QsortCmpBoardClass (name): (%v/5) (%v) expected: (%v)", i, string(theBytes), string(tt.expected[i].Name[:]))
			}

			ReadAt(
				tt.args.shmaddr,
				1200+int(ptttype.BOARD_HEADER_RAW_SZ)*int(the_idx)+int(ptttype.BOARD_HEADER_TITLE_OFFSET),
				ptttype.BOARD_TITLE_SZ,
				unsafe.Pointer(&outTitle),
			)

			theBytes = types.CstrToBytes(outTitle[:])
			if !reflect.DeepEqual(theBytes, tt.expected[i].Title[:]) {
				t.Errorf("QsortCmpBoardClass: (title): (%v/5) %v expected: %v", i, string(theBytes), string(tt.expected[i].Title[:]))
			}
		}
	}
}

func TestMemcmp(t *testing.T) {
	setupTest()
	defer teardownTest()

	shmid, shmaddr, _, _ := CreateShm(testShmKey, 5000, false)
	defer CloseShm(shmid, shmaddr)

	text1 := [13]byte{}
	copy(text1[:], []byte("ABC"))

	WriteAt(
		shmaddr,
		4500,
		13,
		unsafe.Pointer(&text1),
	)

	text2 := [13]byte{}
	copy(text2[:], []byte("456"))

	text3 := [13]byte{}
	copy(text3[:], []byte("C124"))

	text4 := [13]byte{}
	copy(text4[:], []byte("1DE"))

	type args struct {
		shmaddr unsafe.Pointer
		offset  int
		size    uintptr
		cmpaddr unsafe.Pointer
	}
	tests := []struct {
		name     string
		args     args
		expected int
	}{
		// TODO: Add test cases.
		{
			name: "ABC",
			args: args{
				shmaddr: shmaddr,
				offset:  4500,
				size:    uintptr(3),
				cmpaddr: unsafe.Pointer(&text1),
			},
			expected: 0, //ABC: equal
		},
		{
			name: "456",
			args: args{
				shmaddr: shmaddr,
				offset:  4500,
				size:    uintptr(3),
				cmpaddr: unsafe.Pointer(&text2),
			},
			expected: 1, //456: gt
		},
		{
			name: "C124",
			args: args{
				shmaddr: shmaddr,
				offset:  4500,
				size:    uintptr(3),
				cmpaddr: unsafe.Pointer(&text3),
			},
			expected: -1, //C124 lt
		},
		{
			name: "1DE",
			args: args{
				shmaddr: shmaddr,
				offset:  4500,
				size:    uintptr(3),
				cmpaddr: unsafe.Pointer(&text4),
			},
			expected: 1, //1DE gt
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Memcmp(tt.args.shmaddr, tt.args.offset, tt.args.size, tt.args.cmpaddr); got*tt.expected < 0 {
				t.Errorf("Memcmp() = %v, want %v", got, tt.expected)
			}
		})
	}
}
