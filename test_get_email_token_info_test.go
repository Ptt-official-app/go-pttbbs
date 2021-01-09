package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/sirupsen/logrus"
)

func Test_GetEmailTokenInfo(t *testing.T) {
	setupTest()
	defer teardownTest()

	theJwt, err := api.CreateEmailToken("SYSOP", "", "test@ptt.test", api.CONTEXT_CHANGE_EMAIL)
	logrus.Infof("Test_GetEmailTokenInfo: after CreateEmailToken: theJwt: %v e: %v", theJwt, err)
	params0 := &api.GetEmailTokenInfoParams{Jwt: theJwt, Context: api.CONTEXT_CHANGE_EMAIL}

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
				path:     "/emailtoken/info",
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
