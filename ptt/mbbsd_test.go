package ptt

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestLoginQuery(t *testing.T) {
	setupTest()
	defer teardownTest()

	userid1 := ptttype.UserID_t{}
	copy(userid1[:], []byte("SYSOP"))

	type args struct {
		userID *ptttype.UserID_t
		passwd []byte
		ip     *ptttype.IPv4_t
	}
	tests := []struct {
		name     string
		args     args
		expected *ptttype.UserecRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{userID: &userid1, passwd: []byte("123123")},
			expected: testUserecRaw1,
		},
		{
			args:     args{userID: &userid1, passwd: []byte("124")},
			expected: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoginQuery(tt.args.userID, tt.args.passwd, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("LoginQuery() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
