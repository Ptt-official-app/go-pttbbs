package cache

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	log "github.com/sirupsen/logrus"
)

func TestStatInc(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, false, true)
	if err != nil {
		log.Errorf("TestStatInc: unable to NewSHM: e: %v", err)
		return
	}
	defer CloseSHM()

	type args struct {
		stats ptttype.Stat
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{stats: ptttype.STAT_BOARDREC},
			expected: 1,
		},
		{
			args:     args{stats: ptttype.STAT_BOARDREC},
			expected: 2,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := StatInc(tt.args.stats); (err != nil) != tt.wantErr {
				t.Errorf("StatInc() error = %v, wantErr %v", err, tt.wantErr)
			}

			out := Shm.Shm.Statistic[ptttype.STAT_BOARDREC]

			if !reflect.DeepEqual(out, tt.expected) {
				t.Errorf("StatInc() out: %v expected: %v", out, tt.expected)
			}
		})
		wg.Wait()
	}
}

func TestReadStat(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, false, true)
	if err != nil {
		log.Errorf("TestReadStat: unable to NewSHM: e: %v", err)
		return
	}
	defer CloseSHM()

	CleanStat()

	_ = StatInc(ptttype.STAT_BOARDREC)
	_ = StatInc(ptttype.STAT_BOARDREC)
	_ = StatInc(ptttype.STAT_GAMBLE)

	type args struct {
		stats ptttype.Stat
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{ptttype.STAT_BOARDREC},
			expected: 2,
		},
		{
			args:     args{ptttype.STAT_BOARDREC},
			expected: 2,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := ReadStat(tt.args.stats)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadStat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ReadStat() = %v, expected %v", got, tt.expected)
			}
		})
		wg.Wait()
	}
}

func TestCleanStat(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, false, true)
	if err != nil {
		log.Errorf("CleanStat: unable to new shm: e: %v", err)
		return
	}
	defer CloseSHM()

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			CleanStat()

			statsBoardRec, _ := ReadStat(ptttype.STAT_BOARDREC)
			if statsBoardRec != 0 {
				t.Errorf("CleanStat() statsBoardRec: %v", statsBoardRec)
			}
		})
		wg.Wait()
	}
}
