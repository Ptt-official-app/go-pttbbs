package cmsys

import (
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestPttLock(t *testing.T) {
	file, err := os.OpenFile("./testcase/.exist", syscall.O_WRONLY, 0644)
	logrus.Infof("TestPttLock: after open: e: %v", err)
	defer file.Close()

	file2, err := os.OpenFile("./testcase/.exist2", syscall.O_RDWR, 0644)
	logrus.Infof("TestPttLock: after open: e: %v", err)
	defer file.Close()

	type args struct {
		file    *os.File
		offset  int64
		theSize uintptr
		mode    int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{file: file, offset: 10, theSize: 20, mode: syscall.F_WRLCK},
		},
		{
			args: args{file: file, offset: 10, theSize: 20, mode: syscall.F_UNLCK},
		},
		{
			args: args{file: file, offset: 200, theSize: 300, mode: syscall.F_WRLCK},
		},
		{
			args: args{file: file, offset: 200, theSize: 300, mode: syscall.F_UNLCK},
		},
		{
			args: args{file: file2, offset: 10, theSize: 20, mode: syscall.F_WRLCK},
		},
		{
			args: args{file: file2, offset: 10, theSize: 20, mode: syscall.F_UNLCK},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := PttLock(tt.args.file, tt.args.offset, tt.args.theSize, tt.args.mode); (err != nil) != tt.wantErr {
				t.Errorf("PttLock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}

func TestPttLock2(t *testing.T) {
	file, err := os.OpenFile("./testcase/.exist3", syscall.O_WRONLY, 0644)
	logrus.Infof("TestPttLock: after open: e: %v", err)
	defer file.Close()

	type args struct {
		file    *os.File
		offset  int64
		theSize uintptr
		mode    int
	}
	tests := []struct {
		name      string
		args      args
		presleep  int
		postsleep int
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name:      "0th (no blocked)",
			args:      args{file: file, offset: 10, theSize: 20},
			presleep:  0,
			postsleep: 5,
		},
		{
			name:      "1st (blocked, no effect on same pid)",
			args:      args{file: file, offset: 10, theSize: 20},
			presleep:  1,
			postsleep: 1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			postsleep := tt.postsleep
			defer wg.Done()
			logrus.Infof("PttLock2: start: name: %v presleep: %v postsleep: %v", tt.name, tt.presleep, tt.postsleep)
			time.Sleep(time.Duration(tt.presleep) * time.Second)
			file.Seek(0, 0)
			err := PttLock(tt.args.file, tt.args.offset, tt.args.theSize, syscall.F_WRLCK)
			logrus.Infof("PttLock2: after: PttLock: e: %v", err)
			if (err != nil) != tt.wantErr {
				t.Errorf("PttLock() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				logrus.Infof("PttLock2: to sleep2: name: %v postsleep: %v", tt.name, postsleep)
				time.Sleep(time.Duration(postsleep) * time.Second)

				file.Seek(0, 0)
				err := PttLock(tt.args.file, tt.args.offset, tt.args.theSize, syscall.F_UNLCK)
				if err != nil {
					t.Errorf("PttLock() error = %v", err)
				}
			}()

			logrus.Infof("PttLock2: done: name: %v", tt.name)

		})
	}
	wg.Wait()

	logrus.Infof("PttLock2: done")
}

func TestGoPttLock(t *testing.T) {
	filename := "./testcase/.exist3"
	file, err := os.OpenFile(filename, syscall.O_WRONLY, 0644)
	logrus.Infof("TestGoPttLock: after open: e: %v", err)
	defer file.Close()

	type args struct {
		file     *os.File
		filename string
		offset   int64
		theSize  uintptr
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{file: file, filename: filename, offset: 0, theSize: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GoPttLock(tt.args.file, tt.args.filename, tt.args.offset, tt.args.theSize); (err != nil) != tt.wantErr {
				t.Errorf("GoPttLock() error = %v, wantErr %v", err, tt.wantErr)
			}

			GoPttUnlock(tt.args.file, tt.args.filename, tt.args.offset, tt.args.theSize)
		})
	}
}

func TestGoFlock(t *testing.T) {
	filename := "./testcase/.exist3"
	file, err := os.OpenFile(filename, syscall.O_WRONLY, 0644)
	logrus.Infof("TestGoPttLock: after open: e: %v", err)
	defer file.Close()

	fd := file.Fd()

	type args struct {
		fd       uintptr
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{fd: fd, filename: filename},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GoFlock(tt.args.fd, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("GoFlock() error = %v, wantErr %v", err, tt.wantErr)
			}
			GoFunlock(tt.args.fd, tt.args.filename)
		})
	}
}
