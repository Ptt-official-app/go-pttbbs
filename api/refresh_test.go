package api

import (
	"net/http"
	"reflect"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func TestRefresh(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	jwt, _ := CreateToken("SYSOP", "")
	refreshJwt, _ := CreateRefreshToken("SYSOP", "")

	logrus.Infof("TestRefresh: jwt: %v refreshJwt: %v", jwt, refreshJwt)

	params0 := &RefreshParams{
		ClientInfo: "",
		Refresh:    refreshJwt,
	}

	req, _ := http.NewRequest("POST", "http://localhost/refresh", nil)
	req.Header = map[string][]string{
		"Authorization": {"bearer " + jwt},
	}
	c0 := &gin.Context{}
	c0.Request = req

	type args struct {
		remoteAddr string
		params     interface{}
		c          *gin.Context
	}
	tests := []struct {
		name     string
		args     args
		expected interface{}
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{params: params0, c: c0},
			expected: &RefreshResult{UserID: "SYSOP", TokenType: "bearer"},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := Refresh(tt.args.remoteAddr, tt.args.params, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Refresh() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			result, _ := gotResult.(*RefreshResult)
			result.Jwt = ""
			result.Refresh = ""
			if !reflect.DeepEqual(gotResult, tt.expected) {
				t.Errorf("Refresh() = %v, want %v", gotResult, tt.expected)
			}
		})
		wg.Wait()
	}
}
