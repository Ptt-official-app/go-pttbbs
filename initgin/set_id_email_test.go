package initgin

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
)

func Test_SetIDEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	jwt, _ := api.CreateEmailToken("SYSOP", "", "test@ptt.test", api.CONTEXT_SET_ID_EMAIL)

	params0 := &api.SetIDEmailParams{
		Jwt:   jwt,
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
				path:     "/user/SYSOP/setidemail",
				username: "SYSOP",
				passwd:   "123123",
				params:   params0,
			},
		},
	}

	router, _ := InitGin()
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			jwt := getJwt(router, tt.args.username, tt.args.passwd)
			w := httptest.NewRecorder()
			req := setRequest(tt.args.path, tt.args.params, jwt, nil, "POST")
			router.ServeHTTP(w, req)

			body, _ := ioutil.ReadAll(w.Body)
			if w.Code != http.StatusOK {
				t.Errorf("code: %v body: %v", w.Code, string(body))
			}
		})
		wg.Wait()
	}
}
