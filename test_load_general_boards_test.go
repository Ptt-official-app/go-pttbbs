package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
)

func Test_LoadGeneralBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	params := &api.LoadGeneralBoardsParams{
		StartIdx: "",
		NBoards:  4,
	}

	type args struct {
		path   string
		userID string
		passwd string
		params interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			args: args{
				path:   api.LOAD_GENERAL_BOARDS_R,
				userID: "SYSOP",
				passwd: "123123",
				params: params,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := initGin()

			jwt := getJwt(router, tt.args.userID, tt.args.passwd)
			w := httptest.NewRecorder()
			req := setRequest(tt.args.path, params, jwt, nil, "GET")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}
		})
	}

}
