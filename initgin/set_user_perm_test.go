package initgin

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func Test_SetUserPerm(t *testing.T) {
	setupTest()
	defer teardownTest()
	origPerm := testUserec.Userlevel
	newPerm := ptttype.PERM_DEFAULT | ptttype.PERM_ADMIN | ptttype.PERM_LOGINOK

	params0 := &api.SetUserPermParams{
		Perm: newPerm,
	}

	params1 := &api.SetUserPermParams{
		Perm: origPerm,
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
		// TODO: Add test cases.
		{
			args: args{
				path:     "/admin/user/SYSOP/setperm",
				username: "SYSOP",
				passwd:   "123123",
				params:   params0,
			},
		},
		{
			args: args{
				path:     "/admin/user/SYSOP/setperm",
				username: "SYSOP",
				passwd:   "123123",
				params:   params1,
			},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			jwt := getJwt(router, tt.args.username, tt.args.passwd)
			w := httptest.NewRecorder()
			req := setRequest(tt.args.path, tt.args.params, jwt, nil, "POST")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}
		})
		wg.Wait()
	}
}
