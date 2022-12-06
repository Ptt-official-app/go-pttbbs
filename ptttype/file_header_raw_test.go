package ptttype

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func TestFileRefer_Ref(t *testing.T) {
	setupTest()
	defer teardownTest()

	tests := []struct {
		name     string
		f        FileRefer
		expected uint32
	}{
		// TODO: Add test cases.
		{
			f:        0x88888888,
			expected: 0x08888888,
		},
		{
			f:        0xa8888888,
			expected: 0x28888888,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Ref(); got != tt.expected {
				t.Errorf("FileRefer.Ref() = %x, want %x", got, tt.expected)
			}
		})
	}
}

func TestFileRefer_Flag(t *testing.T) {
	setupTest()
	defer teardownTest()

	tests := []struct {
		name     string
		f        FileRefer
		expected uint8
	}{
		// TODO: Add test cases.
		{
			f:        0x88888888,
			expected: 0x1,
		},
		{
			f:        0xa8888888,
			expected: 0x1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := tt.f.Flag(); got != tt.expected {
				t.Errorf("FileRefer.Flag() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func TestFileHeaderRaw_Money(t *testing.T) {
	setupTest()
	defer teardownTest()

	file, _ := os.Open("./testcase/DIR")
	defer file.Close()

	bin, _ := io.ReadAll(file)
	buf := bytes.NewReader(bin)
	records := [1]FileHeaderRaw{}
	_ = types.BinaryRead(buf, binary.LittleEndian, &records)

	type fields struct {
		Filename  Filename_t
		Modified  types.Time4
		Pad       byte
		Recommend int8
		Owner     Owner_t
		Date      Date_t
		Title     Title_t
		Pad2      byte
		Multi     [4]byte
		Filemode  FileMode
		Pad3      [3]byte
	}
	tests := []struct {
		name          string
		record        *FileHeaderRaw
		expectedMoney int32
	}{
		// TODO: Add test cases.
		{
			record:        &records[0],
			expectedMoney: 5,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			f := tt.record
			logrus.Infof("tt.records: %v", tt.record)
			if gotMoney := f.Money(); gotMoney != tt.expectedMoney {
				t.Errorf("FileHeaderRaw.Money() = %v, want %v", gotMoney, tt.expectedMoney)
			}
		})
	}
	wg.Wait()
}

func TestFileHeaderRaw_AnonUID(t *testing.T) {
	setupTest()
	defer teardownTest()

	file, _ := os.Open("./testcase/DIR")
	defer file.Close()

	bin, _ := io.ReadAll(file)
	buf := bytes.NewReader(bin)
	records := [1]FileHeaderRaw{}
	_ = types.BinaryRead(buf, binary.LittleEndian, &records)

	type fields struct {
		Filename  Filename_t
		Modified  types.Time4
		Pad       byte
		Recommend int8
		Owner     Owner_t
		Date      Date_t
		Title     Title_t
		Pad2      byte
		Multi     [4]byte
		Filemode  FileMode
		Pad3      [3]byte
	}
	tests := []struct {
		name            string
		fields          fields
		record          *FileHeaderRaw
		expectedAnonUID int32
	}{
		// TODO: Add test cases.
		{
			record:          &records[0],
			expectedAnonUID: 5,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			f := tt.record
			if gotAnonUID := f.AnonUID(); gotAnonUID != tt.expectedAnonUID {
				t.Errorf("FileHeaderRaw.AnonUID() = %v, want %v", gotAnonUID, tt.expectedAnonUID)
			}
		})
	}
	wg.Wait()
}

func TestFileHeaderRaw_VoteLimits(t *testing.T) {
	setupTest()
	defer teardownTest()

	file, _ := os.Open("./testcase/DIR")
	defer file.Close()

	bin, _ := io.ReadAll(file)
	buf := bytes.NewReader(bin)
	records := [1]FileHeaderRaw{}
	_ = types.BinaryRead(buf, binary.LittleEndian, &records)

	type fields struct {
		Filename  Filename_t
		Modified  types.Time4
		Pad       byte
		Recommend int8
		Owner     Owner_t
		Date      Date_t
		Title     Title_t
		Pad2      byte
		Multi     [4]byte
		Filemode  FileMode
		Pad3      [3]byte
	}
	tests := []struct {
		name     string
		fields   fields
		record   *FileHeaderRaw
		expected *VoteLimits
	}{
		// TODO: Add test cases.
		{
			record:   &records[0],
			expected: &VoteLimits{Post: 5},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			f := tt.record
			if got := f.VoteLimits(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("FileHeaderRaw.VoteLimits() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func TestFileHeaderRaw_VoteLimitPosts(t *testing.T) {
	setupTest()
	defer teardownTest()

	file, _ := os.Open("./testcase/DIR")
	defer file.Close()

	bin, _ := io.ReadAll(file)
	buf := bytes.NewReader(bin)
	records := [1]FileHeaderRaw{}
	_ = types.BinaryRead(buf, binary.LittleEndian, &records)

	type fields struct {
		Filename  Filename_t
		Modified  types.Time4
		Pad       byte
		Recommend int8
		Owner     Owner_t
		Date      Date_t
		Title     Title_t
		Pad2      byte
		Multi     [4]byte
		Filemode  FileMode
		Pad3      [3]byte
	}
	tests := []struct {
		name     string
		fields   fields
		record   *FileHeaderRaw
		expected uint8
	}{
		// TODO: Add test cases.
		{
			record:   &records[0],
			expected: 5,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			f := tt.record
			if got := f.VoteLimitPosts(); got != tt.expected {
				t.Errorf("FileHeaderRaw.VoteLimitPosts() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func TestFileHeaderRaw_VoteLimitLogins(t *testing.T) {
	setupTest()
	defer teardownTest()

	file, _ := os.Open("./testcase/DIR")
	defer file.Close()

	bin, _ := io.ReadAll(file)
	buf := bytes.NewReader(bin)
	records := [1]FileHeaderRaw{}
	_ = types.BinaryRead(buf, binary.LittleEndian, &records)

	type fields struct {
		Filename  Filename_t
		Modified  types.Time4
		Pad       byte
		Recommend int8
		Owner     Owner_t
		Date      Date_t
		Title     Title_t
		Pad2      byte
		Multi     [4]byte
		Filemode  FileMode
		Pad3      [3]byte
	}
	tests := []struct {
		name     string
		fields   fields
		record   *FileHeaderRaw
		expected uint8
	}{
		// TODO: Add test cases.
		{
			record:   &records[0],
			expected: 0,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			f := tt.record
			if got := f.VoteLimitLogins(); got != tt.expected {
				t.Errorf("FileHeaderRaw.VoteLimitLogins() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func TestFileHeaderRaw_VoteLimitRegTime(t *testing.T) {
	setupTest()
	defer teardownTest()

	file, _ := os.Open("./testcase/DIR")
	defer file.Close()

	bin, _ := io.ReadAll(file)
	buf := bytes.NewReader(bin)
	records := [1]FileHeaderRaw{}
	_ = types.BinaryRead(buf, binary.LittleEndian, &records)

	type fields struct {
		Filename  Filename_t
		Modified  types.Time4
		Pad       byte
		Recommend int8
		Owner     Owner_t
		Date      Date_t
		Title     Title_t
		Pad2      byte
		Multi     [4]byte
		Filemode  FileMode
		Pad3      [3]byte
	}
	tests := []struct {
		name     string
		fields   fields
		record   *FileHeaderRaw
		expected uint8
	}{
		// TODO: Add test cases.
		{
			record:   &records[0],
			expected: 0,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			f := &FileHeaderRaw{
				Filename:  tt.fields.Filename,
				Modified:  tt.fields.Modified,
				Pad:       tt.fields.Pad,
				Recommend: tt.fields.Recommend,
				Owner:     tt.fields.Owner,
				Date:      tt.fields.Date,
				Title:     tt.fields.Title,
				Pad2:      tt.fields.Pad2,
				Multi:     tt.fields.Multi,
				Filemode:  tt.fields.Filemode,
				Pad3:      tt.fields.Pad3,
			}
			if got := f.VoteLimitRegTime(); got != tt.expected {
				t.Errorf("FileHeaderRaw.VoteLimitRegTime() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func TestFileHeaderRaw_VoteLimitBadpost(t *testing.T) {
	setupTest()
	defer teardownTest()

	file, _ := os.Open("./testcase/DIR")
	defer file.Close()

	bin, _ := io.ReadAll(file)
	buf := bytes.NewReader(bin)
	records := [1]FileHeaderRaw{}
	_ = types.BinaryRead(buf, binary.LittleEndian, &records)

	type fields struct {
		Filename  Filename_t
		Modified  types.Time4
		Pad       byte
		Recommend int8
		Owner     Owner_t
		Date      Date_t
		Title     Title_t
		Pad2      byte
		Multi     [4]byte
		Filemode  FileMode
		Pad3      [3]byte
	}
	tests := []struct {
		name     string
		fields   fields
		record   *FileHeaderRaw
		expected uint8
	}{
		// TODO: Add test cases.
		{
			record:   &records[0],
			expected: 0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			f := tt.record
			if got := f.VoteLimitBadpost(); got != tt.expected {
				t.Errorf("FileHeaderRaw.VoteLimitBadpost() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func TestFileHeaderRaw_IsDeleted(t *testing.T) {
	setupTest()
	defer teardownTest()

	filename0 := Filename_t{}
	copy(filename0[:], "M.1234567890.A.BCD")

	filename1 := Filename_t{}
	copy(filename1[:], ".deleted")

	owner2 := Owner_t{}
	copy(owner2[:], "-")

	type fields struct {
		Filename  Filename_t
		Modified  types.Time4
		Pad       byte
		Recommend int8
		Owner     Owner_t
		Date      Date_t
		Title     Title_t
		Pad2      byte
		Multi     [4]byte
		Filemode  FileMode
		Pad3      [3]byte
	}
	tests := []struct {
		name     string
		fields   fields
		expected bool
	}{
		// TODO: Add test cases.
		{
			fields:   fields{Filename: filename0},
			expected: false,
		},
		{
			fields:   fields{Filename: filename1},
			expected: true,
		},
		{
			fields:   fields{Owner: owner2},
			expected: true,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			f := &FileHeaderRaw{
				Filename:  tt.fields.Filename,
				Modified:  tt.fields.Modified,
				Pad:       tt.fields.Pad,
				Recommend: tt.fields.Recommend,
				Owner:     tt.fields.Owner,
				Date:      tt.fields.Date,
				Title:     tt.fields.Title,
				Pad2:      tt.fields.Pad2,
				Multi:     tt.fields.Multi,
				Filemode:  tt.fields.Filemode,
				Pad3:      tt.fields.Pad3,
			}
			if got := f.IsDeleted(); got != tt.expected {
				t.Errorf("FileHeaderRaw.IsDeleted() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}
