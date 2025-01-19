package path

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestSetHomePath(t *testing.T) {
	setupTest()
	defer teardownTest()

	userID1 := &ptttype.UserID_t{0x53, 0x59, 0x53, 0x4f, 0x50} // SYSOP

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			args: args{userID: userID1},
			want: "/home/bbs/home/S/SYSOP",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetHomePath(tt.args.userID); got != tt.want {
				t.Errorf("SetHomePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetBPath(t *testing.T) {
	setupTest()
	defer teardownTest()

	boardID1 := &ptttype.BoardID_t{'W', 'h', 'o', 'A', 'm', 'I', 0x00, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e}
	type args struct {
		boardID *ptttype.BoardID_t
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			args: args{boardID: boardID1},
			want: "/home/bbs/boards/W/WhoAmI",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetBPath(tt.args.boardID); got != tt.want {
				t.Errorf("SetBPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetBPath2(t *testing.T) {
	setupTest()
	defer teardownTest()

	IS_DUP_BOARD_DIR_PATH = true
	defer func() {
		IS_DUP_BOARD_DIR_PATH = false
	}()

	boardID1 := &ptttype.BoardID_t{'W', 'h', 'o', 'A', 'm', 'I', 0x00, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e}
	type args struct {
		boardID *ptttype.BoardID_t
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			args: args{boardID: boardID1},
			want: "/home/bbs/boards/W/WhoAmI/WhoAmI",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetBPath(tt.args.boardID); got != tt.want {
				t.Errorf("SetBPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetBFile(t *testing.T) {
	setupTest()
	defer teardownTest()

	boardID1 := &ptttype.BoardID_t{'W', 'h', 'o', 'A', 'm', 'I', 0x00, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e}
	filename1 := "M.1234567890.ABC"

	type args struct {
		boardID  *ptttype.BoardID_t
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{boardID: boardID1, filename: filename1},
			want: "/home/bbs/boards/W/WhoAmI/M.1234567890.ABC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetBFile(tt.args.boardID, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetBFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetBFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetBFile2(t *testing.T) {
	setupTest()
	defer teardownTest()

	IS_DUP_BOARD_DIR_PATH = true
	defer func() {
		IS_DUP_BOARD_DIR_PATH = false
	}()

	boardID1 := &ptttype.BoardID_t{'W', 'h', 'o', 'A', 'm', 'I', 0x00, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e}
	filename1 := "M.1234567890.ABC"

	type args struct {
		boardID  *ptttype.BoardID_t
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{boardID: boardID1, filename: filename1},
			want: "/home/bbs/boards/W/WhoAmI/WhoAmI/M.1234567890.ABC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetBFile(tt.args.boardID, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetBFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetBFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetBNFile(t *testing.T) {
	setupTest()
	defer teardownTest()

	IS_DUP_BOARD_DIR_PATH = true
	defer func() {
		IS_DUP_BOARD_DIR_PATH = false
	}()

	boardID1 := &ptttype.BoardID_t{'W', 'h', 'o', 'A', 'm', 'I', 0x00, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e}
	filename1 := "M.1234567890.ABC"

	type args struct {
		boardID  *ptttype.BoardID_t
		filename string
		idx      int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{boardID: boardID1, filename: filename1, idx: 1},
			want: "/home/bbs/boards/W/WhoAmI/WhoAmI/M.1234567890.ABC.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetBNFile(tt.args.boardID, tt.args.filename, tt.args.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetBNFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetBNFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetBNFile2(t *testing.T) {
	setupTest()
	defer teardownTest()

	boardID1 := &ptttype.BoardID_t{'W', 'h', 'o', 'A', 'm', 'I', 0x00, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e, 0x2e}
	filename1 := "M.1234567890.ABC"

	type args struct {
		boardID  *ptttype.BoardID_t
		filename string
		idx      int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{boardID: boardID1, filename: filename1, idx: 1},
			want: "/home/bbs/boards/W/WhoAmI/M.1234567890.ABC.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetBNFile(tt.args.boardID, tt.args.filename, tt.args.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetBNFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetBNFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDIRPath(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		dirFilename string
		filename    string
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		// TODO: Add test cases.
		{
			args:     args{dirFilename: "boards/W/WhoAmI/.DIR", filename: "test"},
			expected: "boards/W/WhoAmI/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetDIRPath(tt.args.dirFilename, tt.args.filename); got != tt.expected {
				t.Errorf("SetDIRPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}
