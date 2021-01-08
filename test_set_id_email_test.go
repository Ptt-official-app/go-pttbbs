package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
)

func Test_SetIDEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	//until Fri Nov 14 01:28:37 EST 2245
	params0 := &api.SetIDEmailParams{
		Jwt:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGkiOiJ0ZXN0X2NsaWVudGluZm8iLCJlbWwiOiJ0ZXN0QHB0dC50ZXN0IiwiZXhwIjo4NzA1NjAwOTE3LCJzdWIiOiJTWVNPUCJ9.IPa_yF9BYbnnrjwguwPnE7fxpV13bhfgex74-ok-VzE",
		IsSet: true,
	}

	type args struct {
		path     string
		username string
		passwd   string
		params   interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				path:     "/users/SYSOP/setidemail",
				username: "SYSOP",
				passwd:   "123123",
				params:   params0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := initGin()

			jwt := getJwt(router, tt.args.username, tt.args.passwd)
			w := httptest.NewRecorder()
			req := setRequest(tt.args.path, tt.args.params, jwt, nil, "POST")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}
		})
	}

}
