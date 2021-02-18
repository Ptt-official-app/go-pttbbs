package cmbbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/sirupsen/logrus"
)

func TestStampfile(t *testing.T) {
	setupTest()
	defer teardownTest()

	boardFilename := "./testcase/boards/W/WhoAmI"

	filename0 := ptttype.Filename_t{}
	copy(filename0[:], []byte("M.1234567890.ABC"))

	owner0 := ptttype.Owner_t{}
	copy(owner0[:], []byte("SYSOP"))

	title0 := ptttype.Title_t{}
	copy(title0[:], []byte("Re: test0"))

	header0 := &ptttype.FileHeaderRaw{
		Filename: filename0,
		Owner:    owner0,
		Title:    title0,
	}

	expected0 := &ptttype.FileHeaderRaw{}

	type args struct {
		boardFilename string
		header        *ptttype.FileHeaderRaw
	}
	tests := []struct {
		name     string
		args     args
		expected *ptttype.FileHeaderRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{boardFilename: boardFilename, header: header0},
			expected: expected0,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			logrus.Infof("TestSamplefile: heade: %v", tt.args.header)
			_, err := Stampfile(tt.args.boardFilename, tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stampfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.args.header.Filename = expected0.Filename
			tt.args.header.Date = expected0.Date

			testutil.TDeepEqual(t, "got", tt.args.header, expected0)
		})
		wg.Wait()
	}
}

func TestStampfileU(t *testing.T) {
	setupTest()
	defer teardownTest()

	boardFilename := "./testcase/boards/W/WhoAmI"

	filename0 := ptttype.Filename_t{}
	copy(filename0[:], []byte("M.1234567890.ABC"))

	owner0 := ptttype.Owner_t{}
	copy(owner0[:], []byte("SYSOP"))

	title0 := ptttype.Title_t{}
	copy(title0[:], []byte("Re: test0"))

	header0 := &ptttype.FileHeaderRaw{
		Filename: filename0,
		Owner:    owner0,
		Title:    title0,
	}

	expected0 := &ptttype.FileHeaderRaw{
		Owner: owner0,
		Title: title0,
	}

	type args struct {
		boardFilename string
		header        *ptttype.FileHeaderRaw
	}
	tests := []struct {
		name     string
		args     args
		expected *ptttype.FileHeaderRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{boardFilename: boardFilename, header: header0},
			expected: expected0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			_, err := StampfileU(tt.args.boardFilename, tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("StampfileU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.args.header.Filename = expected0.Filename
			tt.args.header.Date = expected0.Date
			testutil.TDeepEqual(t, "got", tt.args.header, expected0)

		})
		wg.Wait()
	}
}

func Test_fhdrStamp(t *testing.T) {
	setupTest()
	defer teardownTest()

	boardFilename := "./testcase/boards/W/WhoAmI"

	owner0 := ptttype.Owner_t{}
	copy(owner0[:], []byte("SYSOP"))

	title0 := ptttype.Title_t{}
	copy(title0[:], []byte("Re: test0"))

	header0 := &ptttype.FileHeaderRaw{
		Owner: owner0,
		Title: title0,
	}

	type args struct {
		boardFilename string
		header        *ptttype.FileHeaderRaw
		theType       ptttype.StampType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{boardFilename: boardFilename, header: header0, theType: ptttype.STAMP_DIR},
		},
		{
			args: args{boardFilename: boardFilename, header: header0, theType: ptttype.STAMP_LINK},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			_, err := fhdrStamp(tt.args.boardFilename, tt.args.header, tt.args.theType)
			if (err != nil) != tt.wantErr {
				t.Errorf("fhdrStamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
		wg.Wait()
	}
}
