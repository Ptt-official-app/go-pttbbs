package cache

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	log "github.com/sirupsen/logrus"
)

func TestSetUMoney(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = LoadUHash()

	userID1 := &ptttype.UserID_t{}
	copy(userID1[:], []byte("SYSOP"))
	SetUserID(1, userID1)

	money := MoneyOf(1)
	log.Infof("TestSetUMoney: MoneyOf(1): money: %v", money)

	type args struct {
		uid   ptttype.Uid
		money int32
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{1, 100},
			want: 100,
		},
		{
			args: args{1, 10000},
			want: 10000,
		},
		{
			args: args{1, 0},
			want: 0,
		},
		{
			args: args{1, money},
			want: money,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := SetUMoney(tt.args.uid, tt.args.money)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetUMoney() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetUMoney() = %v, want %v", got, tt.want)
			}
		})
		wg.Wait()
	}
}

func TestDeUMoney(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = LoadUHash()

	money := MoneyOf(1)
	log.Infof("TestDeUMoney: MoneyOf(1): money: %v", money)

	defer SetUMoney(1, 0)

	type args struct {
		uid   ptttype.Uid
		money int32
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{1, 100},
			want: 100,
		},
		{
			args: args{1, -50},
			want: 50,
		},
		{
			args: args{1, -200},
			want: 0,
		},
		{
			args: args{1, 100},
			want: 100,
		},
		{
			args: args{1, 300},
			want: 400,
		},
		{
			args: args{1, -150},
			want: 250,
		},
		{
			args: args{1, -250},
			want: 0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := DeUMoney(tt.args.uid, tt.args.money)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeUMoney() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeUMoney() = %v, want %v", got, tt.want)
			}
		})
		wg.Wait()
	}
}
