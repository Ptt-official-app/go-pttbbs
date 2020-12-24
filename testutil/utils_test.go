package testutil

import (
	"testing"
)

func TestTDeepEqual(t *testing.T) {
	type temp struct {
		Temp4 int
		temp3 []*temp
	}
	type temp2 struct {
		Temp3 []*temp
	}

	got0 := []*temp2{
		{
			Temp3: []*temp{{}, {}, {}},
		},
		{
			Temp3: nil,
		},
		{
			Temp3: []*temp{},
		},
	}

	got1 := []*temp2{
		{
			Temp3: nil,
		},
		{
			Temp3: []*temp{{}, {}, {}},
		},
		{
			Temp3: []*temp{},
		},
	}
	type args struct {
		t        *testing.T
		prompt   string
		got      interface{}
		expected interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			args: args{
				t:        t,
				prompt:   "test0",
				got:      got0,
				expected: got0,
			},
		},
		{
			args: args{
				t:        t,
				prompt:   "test1",
				got:      got1,
				expected: got1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TDeepEqual(tt.args.t, tt.args.prompt, tt.args.got, tt.args.expected)
		})
	}
}
