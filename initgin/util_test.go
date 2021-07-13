package initgin

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"
	"github.com/sirupsen/logrus"
)

func setRequest(path string, params interface{}, jwt string, headers map[string]string, method string) *http.Request {
	var jsonBytes []byte

	if method == "GET" {
		v, _ := query.Values(params)
		path = path + "?" + v.Encode()
	} else {
		jsonBytes, _ = json.Marshal(params)
	}

	req, _ := http.NewRequest(method, withPrefix(path), bytes.NewBuffer(jsonBytes))

	req.Header.Set("Host", "localhost:5678")
	req.Header.Set("X-Forwarded-For", "127.0.0.1:5679")
	if jwt != "" {
		req.Header.Set("Authorization", "bearer "+jwt)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	return req
}

func getJwt(router *gin.Engine, userID string, passwd string) string {
	w := httptest.NewRecorder()
	loginParams := &api.LoginParams{Username: userID, Passwd: passwd}
	req := setRequest(api.LOGIN_R, loginParams, "", nil, "POST")
	router.ServeHTTP(w, req)

	body, _ := ioutil.ReadAll(w.Body)

	if w.Code != 200 {
		logrus.Errorf("getJwt: unable to get Jwt: userID: %v passwd: %v code: %v body: %v", userID, passwd, w.Code, string(body))
	}

	resultLogin := &api.LoginResult{}
	_ = json.Unmarshal(body, resultLogin)
	jwt := resultLogin.Jwt

	return jwt
}

func Test_setRequest(t *testing.T) {
	params := "test"
	headers := make(map[string]string)
	headers["Contetn-Type"] = "application/json"

	type args struct {
		path    string
		params  interface{}
		jwt     string
		headers map[string]string
	}
	tests := []struct {
		name     string
		args     args
		expected *http.Request
	}{
		// TODO: Add test cases.
		{
			args: args{params: params, headers: nil},
		},
		{
			args: args{params: params, headers: headers},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setRequest(tt.args.path, tt.args.params, tt.args.jwt, tt.args.headers, "POST")
			if got == nil {
				t.Errorf("unable to setRequest")
			}
		})
	}
}
