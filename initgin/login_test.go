package initgin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
	log "github.com/sirupsen/logrus"
)

func Test_Login(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		path   string
		params interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			args: args{
				path: api.LOGIN_R,
				params: &api.LoginParams{
					Username: "SYSOP",
					Passwd:   "123123",
				},
			}, // json: {}
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			router, _ := InitGin()

			w := httptest.NewRecorder()
			req := setRequest(tt.args.path, tt.args.params, "", nil, "POST")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}

			body, _ := ioutil.ReadAll(w.Body)

			results := &api.LoginResult{}
			json.Unmarshal(body, results)

			jwt := results.Jwt

			log.Infof("jwt: %v", jwt)

		})
		wg.Wait()
	}
}
