package fav

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestFavType_copyTo(t *testing.T) {
	ft0 := &FavType{FAVT_BOARD, 1, &FavBoard{1, 0, 0}}
	ft1 := &FavType{}
	type args struct {
		ft1 *FavType
	}
	tests := []struct {
		name string
		ft   *FavType
		args args
	}{
		// TODO: Add test cases.
		{
			ft:   ft0,
			args: args{ft1: ft1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ft := tt.ft
			ft.copyTo(tt.args.ft1)

			testutil.TDeepEqual(t, "ft", tt.args.ft1, ft)
		})
	}
}

func TestFavType_GetID(t *testing.T) {
	ft0 := &FavType{FAVT_BOARD, FAVH_FAV, &FavBoard{2, 0, 0}}
	ft1 := &FavType{FAVT_LINE, FAVH_FAV, &FavLine{3}}
	ft2 := &FavType{FAVT_FOLDER, FAVH_FAV, &FavFolder{Fid: 4}}
	tests := []struct {
		name          string
		ft            *FavType
		expectedTheID int
	}{
		// TODO: Add test cases.
		{
			ft:            ft0,
			expectedTheID: 2,
		},
		{
			ft:            ft1,
			expectedTheID: 3,
		},
		{
			ft:            ft2,
			expectedTheID: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ft := tt.ft
			if gotTheID := ft.GetID(); gotTheID != tt.expectedTheID {
				t.Errorf("FavType.GetID() = %v, want %v", gotTheID, tt.expectedTheID)
			}
		})
	}
}
