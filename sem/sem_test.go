package sem

import (
	"errors"
	"os"
	"reflect"
	"syscall"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_errnoErr(t *testing.T) {
	type args struct {
		errno syscall.Errno
	}
	tests := []struct {
		name     string
		args     args
		expected error
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{syscall.EEXIST},
			expected: os.ErrExist,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errnoErr(tt.args.errno)

			if (err != nil) != tt.wantErr {
				t.Errorf("errnoErr() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, tt.expected) {
				t.Errorf("errnoErr() err: %v expected: %v", err, tt.expected)
			}
		})
	}
}

func TestSemGet(t *testing.T) {
	testSem := &Semaphore{testSemKey, 1}
	type args struct {
		key   int
		nsems int
		flags int
	}
	tests := []struct {
		name     string
		args     args
		expected *Semaphore
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{key: testSemKey, nsems: 1, flags: IPC_CREAT | IPC_EXCL | SEM_A | SEM_R},
			expected: testSem,
			wantErr:  false,
		},
		{
			args:     args{key: testSemKey, nsems: 1, flags: IPC_CREAT | IPC_EXCL | SEM_A | SEM_R},
			expected: testSem,
			wantErr:  true,
		},
		{
			args:     args{key: testSemKey, nsems: 1, flags: SEM_A | SEM_R},
			expected: testSem,
			wantErr:  false,
		},
	}
	var firstGot *Semaphore
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SemGet(tt.args.key, tt.args.nsems, tt.args.flags)
			log.Infof("After SemGet: got :%v e: %v", got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("SemGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				return
			}
			if !reflect.DeepEqual(got.nsems, tt.expected.nsems) {
				t.Errorf("SemGet() = %v, expected %v", got, tt.expected)
			}

			if firstGot == nil {
				firstGot = got
			}

		})
	}
	if firstGot != nil {
		firstGot.Destroy()
	}
}

func TestSemaphore_SetVal(t *testing.T) {
	s, err := SemGet(testSemKey, 1, IPC_CREAT|SEM_A|SEM_R)
	if err != nil {
		return
	}
	defer s.Destroy()

	type fields struct {
		semid int
		nsems int
	}
	type args struct {
		semNum int
		val    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{0, 1},
		},
		{
			args: args{0, 2},
		},
		{
			args: args{0, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.SetVal(tt.args.semNum, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Semaphore.SetVal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			getVal, err := s.GetVal(tt.args.semNum)
			if err != nil {
				t.Errorf("unable to GetVal: e: %v", err)
				return
			}

			if !reflect.DeepEqual(tt.args.val, getVal) {
				t.Errorf("Semaphore.SetVal: val: %v expected: %v", getVal, tt.args.val)
			}
		})
	}
}

func TestSemaphore_Wait(t *testing.T) {
	s, err := SemGet(testSemKey, 1, IPC_CREAT|SEM_A|SEM_R)
	if err != nil {
		return
	}
	defer s.Destroy()

	type fields struct {
		semid int
		nsems int
	}
	type args struct {
		semNum int
	}
	tests := []struct {
		name     string
		fields   fields
		theSet   int
		expected int
		args     args
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{0},
			theSet:   1,
			expected: 0,
		},
		{
			args:     args{0},
			theSet:   2,
			expected: 1,
		},
		{
			args:     args{0},
			theSet:   3,
			expected: 2,
		},
		{
			args:     args{0},
			theSet:   100,
			expected: 99,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.SetVal(tt.args.semNum, tt.theSet)

			if err := s.Wait(tt.args.semNum); (err != nil) != tt.wantErr {
				t.Errorf("Semaphore.Wait() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				defer s.Post(tt.args.semNum)
				val, err := s.GetVal(tt.args.semNum)
				if err != nil {
					return
				}

				assert.Equal(t, val, tt.expected)
			}
		})
	}
}

func TestSemaphore_Post(t *testing.T) {
	s, err := SemGet(testSemKey, 1, IPC_CREAT|SEM_A|SEM_R)
	if err != nil {
		return
	}
	defer s.Destroy()

	type fields struct {
		semid int
		nsems int
	}
	type args struct {
		semNum int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		theSet   int
		expected int
		wantErr  bool
	}{
		// TODO: Add test cases.
		// TODO: Add test cases.
		{
			args:     args{0},
			theSet:   0,
			expected: 1,
		},
		{
			args:     args{0},
			theSet:   1,
			expected: 2,
		},
		{
			args:     args{0},
			theSet:   2,
			expected: 3,
		},
		{
			args:     args{0},
			theSet:   99,
			expected: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.SetVal(tt.args.semNum, tt.theSet)

			if err := s.Post(tt.args.semNum); (err != nil) != tt.wantErr {
				t.Errorf("Semaphore.Post() error = %v, wantErr %v", err, tt.wantErr)
			}

			val, err := s.GetVal(tt.args.semNum)
			if err != nil {
				return
			}

			assert.Equal(t, val, tt.expected)
		})
	}
}
