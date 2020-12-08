package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/gin-gonic/gin"
)

func setRequest(path string, params interface{}, jwt string, headers map[string]string) *http.Request {
	jsonStr, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(jsonStr))

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
	loginParams := &api.LoginParams{UserID: userID, Passwd: passwd}
	req := setRequest(login_r, loginParams, "", nil)
	router.ServeHTTP(w, req)

	body, _ := ioutil.ReadAll(w.Body)
	resultLogin := &api.LoginResult{}
	_ = json.Unmarshal(body, resultLogin)
	jwt := resultLogin.Jwt

	return jwt
}
