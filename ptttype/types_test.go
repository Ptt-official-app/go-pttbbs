package ptttype

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

func TestToSortIdx(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		str string
	}
	tests := []struct {
		name     string
		args     args
		expected SortIdx
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{"sdf"},
			expected: -1,
			wantErr:  true,
		},
		{
			args:     args{"-1"},
			expected: -1,
			wantErr:  false,
		},
		{
			args:     args{"0"},
			expected: 0,
			wantErr:  false,
		},
		{
			args:     args{"231"},
			expected: 231,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToSortIdx(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToSortIdx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ToSortIdx() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSortIdx_String(t *testing.T) {
	setupTest()
	defer teardownTest()

	tests := []struct {
		name     string
		s        SortIdx
		expected string
	}{
		// TODO: Add test cases.
		{
			s:        -1,
			expected: "",
		},
		{
			s:        0,
			expected: "0",
		},
		{
			s:        123,
			expected: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.expected {
				t.Errorf("SortIdx.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBoardTitle_t_RealTitle(t *testing.T) {
	setupTest()
	defer teardownTest()

	boardTitle0 := &BoardTitle_t{}
	copy(boardTitle0[:], types.Utf8ToBig5("CPBL ◎四海之內皆兄弟"))
	expected0 := types.Utf8ToBig5("四海之內皆兄弟")

	boardTitle1 := &BoardTitle_t{}
	copy(boardTitle1[:], types.Utf8ToBig5("*CPBL◎四海之內皆兄弟"))
	expected1 := types.Utf8ToBig5("四海之內皆兄弟")

	boardTitle2 := &BoardTitle_t{}
	copy(boardTitle2[:], types.Utf8ToBig5("里肌 ◎SYSOP"))
	expected2 := types.Utf8ToBig5("SYSOP")

	tests := []struct {
		name     string
		tr       *BoardTitle_t
		expected []byte
	}{
		// TODO: Add test cases.
		{
			tr:       boardTitle0,
			expected: expected0,
		},
		{
			tr:       boardTitle1,
			expected: expected1,
		},
		{
			tr:       boardTitle2,
			expected: expected2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.RealTitle(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("BoardTitle_t.RealTitle() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBoardTitle_t_BoardClass(t *testing.T) {
	setupTest()
	defer teardownTest()

	boardTitle0 := &BoardTitle_t{}
	copy(boardTitle0[:], types.Utf8ToBig5("CPBL ◎四海之內皆兄弟"))
	expected0 := []byte("CPBL")

	boardTitle1 := &BoardTitle_t{}
	copy(boardTitle1[:], types.Utf8ToBig5("*CPBL◎四海之內皆兄弟"))
	expected1 := []byte("*CPBL")

	boardTitle2 := &BoardTitle_t{}
	copy(boardTitle2[:], types.Utf8ToBig5("里肌 ◎SYSOP"))
	expected2 := types.Utf8ToBig5("里肌")

	boardTitle3 := &BoardTitle_t{}
	copy(boardTitle3[:], types.Utf8ToBig5("里CP ◎SYSOP"))
	expected3 := types.Utf8ToBig5("里CP")

	tests := []struct {
		name     string
		tr       *BoardTitle_t
		expected []byte
	}{
		// TODO: Add test cases.
		{
			tr:       boardTitle0,
			expected: expected0,
		},
		{
			tr:       boardTitle1,
			expected: expected1,
		},
		{
			tr:       boardTitle2,
			expected: expected2,
		},
		{
			tr:       boardTitle3,
			expected: expected3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.BoardClass(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("BoardTitle_t.BoardClass() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBoardTitle_t_BoardType(t *testing.T) {
	setupTest()
	defer teardownTest()

	boardTitle0 := &BoardTitle_t{}
	copy(boardTitle0[:], types.Utf8ToBig5("CPBL ◎四海之內皆兄弟"))
	expected0 := types.Utf8ToBig5("◎")

	boardTitle1 := &BoardTitle_t{}
	copy(boardTitle1[:], types.Utf8ToBig5("*CPBL◎四海之內皆兄弟"))
	expected1 := types.Utf8ToBig5("◎")

	boardTitle2 := &BoardTitle_t{}
	copy(boardTitle2[:], types.Utf8ToBig5("里肌 ◎SYSOP"))
	expected2 := types.Utf8ToBig5("◎")

	tests := []struct {
		name     string
		tr       *BoardTitle_t
		expected []byte
	}{
		// TODO: Add test cases.
		{
			tr:       boardTitle0,
			expected: expected0,
		},
		{
			tr:       boardTitle1,
			expected: expected1,
		},
		{
			tr:       boardTitle2,
			expected: expected2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.BoardType(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("BoardTitle_t.BoardType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBM_t_ToBMs(t *testing.T) {
	setupTest()
	defer teardownTest()

	bm0 := &BM_t{}
	copy(bm0[:], []byte("te1/te2/te3\x00te4"))
	expectedbm0 := &UserID_t{}
	copy(expectedbm0[:], []byte("te1"))
	expectedbm1 := &UserID_t{}
	copy(expectedbm1[:], []byte("te2"))
	expectedbm2 := &UserID_t{}
	copy(expectedbm2[:], []byte("te3"))
	expected0 := []*UserID_t{expectedbm0, expectedbm1, expectedbm2}
	tests := []struct {
		name     string
		bm       *BM_t
		expected []*UserID_t
	}{
		// TODO: Add test cases.
		{
			bm:       bm0,
			expected: expected0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bm.ToBMs()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("BM_t.ToBMs() = %v, want %v", got, tt.expected)
			}
			for idx, each := range got {
				if idx >= len(tt.expected) {
					t.Errorf("BM_t: (%v/%v) %v", idx, len(got), each)
					continue
				}
				if !reflect.DeepEqual(each, tt.expected[idx]) {
					t.Errorf("BM_t: (%v/%v) %v want: %v", idx, len(got), each, tt.expected[idx])
				}
			}
		})
	}
}

func TestFilename_t_CreateTime(t *testing.T) {
	setupTest()
	defer teardownTest()

	f0 := &Filename_t{}
	copy(f0[:], []byte("M.1234567890.A.123"))

	expected0 := types.Time4(1234567890)

	tests := []struct {
		name     string
		f        *Filename_t
		expected types.Time4
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			f:        f0,
			expected: expected0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.CreateTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("Filename_t.CreateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Filename_t.CreateTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFilename_t_Postfix(t *testing.T) {
	setupTest()
	defer teardownTest()

	f0 := &Filename_t{}
	copy(f0[:], []byte("M.1234567890.A.123"))

	expected0 := []byte("123")

	tests := []struct {
		name     string
		f        *Filename_t
		expected []byte
	}{
		// TODO: Add test cases.
		{
			f:        f0,
			expected: expected0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Postfix(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Filename_t.Postfix() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFilename_t_ToAidu(t *testing.T) {
	setupTest()
	defer teardownTest()

	f0 := &Filename_t{}
	copy(f0[:], []byte("M.1234567890.A.123"))

	expected0 := Aidu(0x499602d2123)

	f1 := &Filename_t{}
	copy(f1[:], []byte("M.1607937174.A.081"))

	expected1 := Aidu(0x5fd72c96081)

	f2 := &Filename_t{}
	copy(f2[:], []byte("M.1607202239.A.30D"))

	expected2 := Aidu(0x5fcbf5bf30d)

	tests := []struct {
		name     string
		f        *Filename_t
		expected Aidu
	}{
		// TODO: Add test cases.
		{
			f:        f0,
			expected: expected0,
		},
		{
			f:        f1,
			expected: expected1,
		},
		{
			f:        f2,
			expected: expected2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.ToAidu(); got != tt.expected {
				t.Errorf("Filename_t.ToAidu() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFilename_t_Type(t *testing.T) {
	setupTest()
	defer teardownTest()

	f0 := &Filename_t{}
	copy(f0[:], []byte("M.1234567890.A.123"))

	f1 := &Filename_t{}
	copy(f1[:], []byte("G.1234567890.A.123"))

	tests := []struct {
		name     string
		f        *Filename_t
		expected RecordType
	}{
		// TODO: Add test cases.
		{
			f:        f0,
			expected: RECORD_TYPE_M,
		},
		{
			f:        f1,
			expected: RECORD_TYPE_G,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Type(); got != tt.expected {
				t.Errorf("Filename_t.Type() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAidu_Type(t *testing.T) {
	setupTest()
	defer teardownTest()

	a0 := Aidu(0x0499602d2123)
	expected0 := RECORD_TYPE_M
	a1 := Aidu(0x1499602d2123)
	expected1 := RECORD_TYPE_G

	tests := []struct {
		name     string
		a        Aidu
		expected RecordType
	}{
		// TODO: Add test cases.
		{
			a:        a0,
			expected: expected0,
		},
		{
			a:        a1,
			expected: expected1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Type(); got != tt.expected {
				t.Errorf("Aidu.Type() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAidu_Time(t *testing.T) {
	setupTest()
	defer teardownTest()

	a0 := Aidu(0x0499602d2123)
	expected0 := types.Time4(1234567890)

	tests := []struct {
		name     string
		a        Aidu
		expected types.Time4
	}{
		// TODO: Add test cases.
		{
			a:        a0,
			expected: expected0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Time(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Aidu.Time() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAidu_ToFN(t *testing.T) {
	setupTest()
	defer teardownTest()

	a0 := Aidu(0x0499602d2123)
	expected0 := &Filename_t{}
	copy(expected0[:], []byte("M.1234567890.A.123"))

	a1 := Aidu(0x5fcbf5bf30d)
	expected1 := &Filename_t{}
	copy(expected1[:], []byte("M.1607202239.A.30D"))

	tests := []struct {
		name     string
		a        Aidu
		expected *Filename_t
	}{
		// TODO: Add test cases.
		{
			a:        a0,
			expected: expected0,
		},
		{
			a:        a1,
			expected: expected1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.ToFN(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Aidu.ToFN() = %v, want %v", got, tt.expected)
				for idx, each := range got {
					if each != tt.expected[idx] {
						t.Errorf("Aidu: (%v/%v) %v want: %v", idx, len(got), each, tt.expected[idx])
					}
				}
			}
		})
	}
}

func TestAidu_ToAidc(t *testing.T) {
	setupTest()
	defer teardownTest()

	a0 := Aidu(0x0499602d2123)
	expected0 := &Aidc{}
	copy(expected0[:], []byte("19bWBI4Z"))

	f1 := &Filename_t{}
	copy(f1[:], []byte("M.1607937174.A.081"))
	a1 := f1.ToAidu()
	log.Infof("f1: %v a1: %x", f1, a1)
	expected1 := &Aidc{}
	copy(expected1[:], []byte("1VrooM21"))

	a2 := Aidu(0x5fcbf5bf30d)
	expected2 := &Aidc{}
	copy(expected2[:], []byte("1Vo_M_CD"))

	tests := []struct {
		name     string
		a        Aidu
		expected *Aidc
	}{
		// TODO: Add test cases.
		{
			a:        a0,
			expected: expected0,
		},
		{
			a:        a1,
			expected: expected1,
		},
		{
			a:        a2,
			expected: expected2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.ToAidc(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Aidu.ToAidc() = %v, want %v", string(got[:]), string(tt.expected[:]))
			}
		})
	}
}

func TestAidc_ToAidu(t *testing.T) {
	setupTest()
	defer teardownTest()

	a0 := &Aidc{}
	copy(a0[:], []byte("19bWBI4Z"))
	expected0 := Aidu(0x0499602d2123)

	a1 := &Aidc{}
	copy(a1[:], []byte("1VrooM21"))
	expected1 := Aidu(0x5fd72c96081)

	tests := []struct {
		name         string
		a            *Aidc
		expectedAidu Aidu
	}{
		// TODO: Add test cases.
		{
			a:            a0,
			expectedAidu: expected0,
		},
		{
			a:            a1,
			expectedAidu: expected1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAidu := tt.a.ToAidu(); gotAidu != tt.expectedAidu {
				t.Errorf("Aidc.ToAidu() = %x, want %x", gotAidu, tt.expectedAidu)
			}
		})
	}
}

func TestOwner_t_ToUserID(t *testing.T) {
	setupTest()
	defer teardownTest()

	o := &Owner_t{}
	copy(o[:], []byte("test1."))

	u := &UserID_t{}
	copy(u[:], []byte("test1"))
	tests := []struct {
		name     string
		o        *Owner_t
		expected *UserID_t
	}{
		// TODO: Add test cases.
		{
			o:        o,
			expected: u,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.ToUserID(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Owner_t.ToUserID() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUserID_t_IsValid(t *testing.T) {
	setupTest()
	defer teardownTest()

	userID0 := UserID_t{}

	userID1 := UserID_t{}
	copy(userID1[:], []byte("S"))

	userID2 := UserID_t{}
	copy(userID2[:], []byte("SYSOP"))

	userID3 := UserID_t{}
	copy(userID3[:], []byte("S1234567891234"))

	userID4 := UserID_t{}
	copy(userID4[:], []byte("SYSOP,-"))

	userID5 := UserID_t{}
	copy(userID5[:], []byte("SYSOP1"))

	userID6 := UserID_t{}
	copy(userID6[:], []byte("1SYSOP"))

	userID7 := UserID_t{}
	copy(userID7[:], []byte("S1"))

	userID8 := UserID_t{}
	copy(userID8[:], []byte("Ss"))

	tests := []struct {
		name     string
		u        *UserID_t
		expected bool
	}{
		// TODO: Add test cases.
		{
			name:     "nil",
			u:        nil,
			expected: false,
		},
		{
			name:     "",
			u:        &userID0,
			expected: false,
		},
		{
			name:     "S",
			u:        &userID1,
			expected: false,
		},
		{
			name:     "SYSOP",
			u:        &userID2,
			expected: true,
		},
		{
			name:     "too long",
			u:        &userID3,
			expected: false,
		},
		{
			name:     "not alnum",
			u:        &userID4,
			expected: false,
		},
		{
			name:     "SYSOP1",
			u:        &userID5,
			expected: true,
		},
		{
			name:     "1SYSOP (not alpha in 0st char)",
			u:        &userID6,
			expected: false,
		},
		{
			name:     "S1 (alnum)",
			u:        &userID7,
			expected: true,
		},
		{
			name:     "Ss (all alpha)",
			u:        &userID7,
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsValid(); got != tt.expected {
				t.Errorf("UserID_t.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}
