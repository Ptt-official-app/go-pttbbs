package cmbbs

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestSubjectEx(t *testing.T) {
	title0 := &ptttype.Title_t{}
	copy(title0[:], []byte("Re: Re: Fw: test0"))
	expectedType0 := ptttype.SUBJECT_FORWARD
	expected0 := []byte("test0")
	type args struct {
		title *ptttype.Title_t
	}
	tests := []struct {
		name              string
		args              args
		expectedTitleType ptttype.SubjectType
		expectedNewTitle  []byte
	}{
		// TODO: Add test cases.
		{
			args:              args{title: title0},
			expectedTitleType: expectedType0,
			expectedNewTitle:  expected0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitleType, gotNewTitle := SubjectEx(tt.args.title)
			if !reflect.DeepEqual(gotTitleType, tt.expectedTitleType) {
				t.Errorf("SubjectEx() gotTitleType = %v, want %v", gotTitleType, tt.expectedTitleType)
			}
			if !reflect.DeepEqual(gotNewTitle, tt.expectedNewTitle) {
				t.Errorf("SubjectEx() gotNewTitle = %v, want %v", gotNewTitle, tt.expectedNewTitle)
			}
		})
	}
}
