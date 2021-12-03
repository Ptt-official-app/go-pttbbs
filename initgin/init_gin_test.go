package initgin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/sirupsen/logrus"
)

func TestInitGin(t *testing.T) {
	setupTest()
	defer teardownTest()

	// variables
	_ = ptt.SetupNewUser(testNewPostUserRaw1)
	_ = ptt.SetupNewUser(testUserecRaw3)

	jwtEmail, _ := api.CreateEmailToken("SYSOP", "", "test@ptt.test", api.CONTEXT_CHANGE_EMAIL)
	jwtIDEmail, _ := api.CreateEmailToken("SYSOP", "", "test@ptt.test", api.CONTEXT_SET_ID_EMAIL)
	jwtSysop2, _ := api.CreateToken("SYSOP", "")

	class0 := []byte("test")
	title0 := []byte("this is a test")
	content0 := [][]byte{[]byte("test1"), []byte("test2")}
	filename0 := &ptttype.Filename_t{}
	copy(filename0[:], []byte("M.1607202239.A.30D"))
	articleID0 := bbs.ToArticleID(filename0)
	ip0 := "127.0.0.1"

	boardSummary, err := bbs.CreateBoard("test", 5, "mewboard0", []byte("CPBL"), []byte("new-board"), nil, 0, 0, 0, false)
	logrus.Infof("after CreateBoard: mewboard0: boardSummary: %v e: %v", boardSummary, err)

	forwardBoardSummary, err := bbs.CreateBoard("test", 5, "fwboard0", []byte("CPBL"), []byte("fw-board"), nil, 0, 0, 0, false)

	logrus.Infof("after CreateBoard: fwboard0: forwardBoardSummary: %v e: %v", forwardBoardSummary, err)

	articleSummary0, err := bbs.CreateArticle("test", boardSummary.BBoardID, class0, title0, content0, ip0)
	if err != nil {
		t.Errorf("unable to create article: e: %v", err)
	}

	gotContent0, _, oldSum0, _ := bbs.GetArticle("test", boardSummary.BBoardID, articleSummary0.ArticleID, 0, true)

	editContent0 := [][]byte{
		{
			0xa7, 0x40, 0xaa, 0xcc, 0x3a, 0x20, 'A', '1', ' ', // 作者: A1
			0x28, 0xaf, 0xab, //(神
			0x29, 0x20, 0xac, 0xdd, 0xaa, 0x4f, //) 看板
			0x3a, 0x20, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, //: WhoAmI
		},
		{
			0xbc, 0xd0, 0xc3, 0x44, 0x3a, 0x20, 0x5b, 0x74, 0x65, 0x73, // 標題: [tes
			0x74, 0x5d, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, // t] this is
			0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, // a test
		},
		{
			0xae, 0xc9, 0xb6, 0xa1, 0x3a, 0x20, // 時間:
		},
		{},
		{
			0x74, 0x65, 0x73, 0x74, 0x38, // test8
		},
		{
			0x74, 0x65, 0x73, 0x74, 0x39, // test9
		},
		{},
		{
			0x2d, 0x2d, //--
		},
		{
			0xa1, 0xb0, 0x20, 0xb5, 0x6f, 0xab, 0x48, 0xaf, 0xb8, 0x3a, //※ 發信站:
			0x20, 0xb7, 0x73, 0xa7, 0xe5, 0xbd, 0xf0, 0xbd, 0xf0, 0x28, // 新批踢踢(
			0x70, 0x74, 0x74, 0x32, 0x2e, 0x63, 0x63, 0x29, 0x2c, 0x20, // ptt2.cc),
			0xa8, 0xd3, 0xa6, 0xdb, 0x3a, 0x20, 0x31, 0x32, 0x37, 0x2e, // 來自: 127.
			0x30, 0x2e, 0x30, 0x2e, 0x31, // 0.0.1
		},
		{
			0xa1, 0xb0, 0x20, 0xa4, 0xe5, 0xb3, 0xb9, 0xba, 0xf4, //※ 文章網
			0xa7, 0x7d, 0x3a, 0x20, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, // 址: http:/
			0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, // /localhost
			0x2f, 0x62, 0x62, 0x73, 0x2f, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, 0x2f, 0x4d, 0x2e, // /bbs/WhoAmI/M.
		},
	}

	// router
	router, _ := InitGin()
	jwtSysop := getJwt(router, "SYSOP", "123123")
	jwtA1 := getJwt(router, "A1", "123123")
	jwtTest := getJwt(router, "test", "123123")
	logrus.Infof("TestInitGin: jwtSysop: %v jwtA1: %v jwtTest: %v", jwtSysop, jwtA1, jwtTest)

	type args struct {
		path   string
		jwt    string
		params interface{}
		method string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Index",
			args: args{
				path:   api.INDEX_R,
				jwt:    jwtSysop,
				method: "POST",
			},
		},
		{
			name: "AttemptChangeEmail",
			args: args{
				path: "/user/SYSOP/attemptchangeemail",
				jwt:  jwtSysop,
				params: &api.AttemptChangeEmailParams{
					ClientInfo: "test_clientinfo",
					Passwd:     "123123",
					Email:      "test@ptt.test",
				},
				method: "POST",
			},
		},
		{
			name: "AttemptSetIDEmail",
			args: args{
				path: "/user/SYSOP/attemptsetidemail",
				jwt:  jwtSysop,
				params: &api.AttemptSetIDEmailParams{
					ClientInfo: "test_clientinfo",
					Passwd:     "123123",
					Email:      "test@ptt.test",
				},
				method: "POST",
			},
		},
		{
			name: "ChangeEmail",
			args: args{
				path: "/user/SYSOP/changeemail",
				jwt:  jwtSysop,
				params: &api.ChangeEmailParams{
					Jwt: jwtEmail,
				},
				method: "POST",
			},
		},
		{
			name: "ChangePasswd",
			args: args{
				path: "/user/SYSOP/changepasswd",
				jwt:  jwtSysop,
				params: &api.ChangePasswdParams{
					ClientInfo: "test_clientinfo",
					OrigPasswd: "123123",
					Passwd:     "123124",
				},
				method: "POST",
			},
		},
		{
			name: "ChangePasswd-2",
			args: args{
				path: "/user/SYSOP/changepasswd",
				jwt:  jwtSysop,
				params: &api.ChangePasswdParams{
					ClientInfo: "test_clientinfo",
					OrigPasswd: "123124",
					Passwd:     "123123",
				},
				method: "POST",
			},
		},
		{
			name: "CheckExistsUser",
			args: args{
				path: api.CHECK_EXISTS_USER_R,
				jwt:  jwtSysop,
				params: &api.CheckExistsUserParams{
					Username: "SYSOP",
				},
				method: "POST",
			},
		},
		{
			name: "CreateArticle",
			args: args{
				path: "/board/10_WhoAmI/article",
				jwt:  jwtA1,
				params: &api.CreateArticleParams{
					PostType: class0,
					Title:    title0,
					Content:  content0,
				},
				method: "POST",
			},
		},
		{
			name: "CreateBoard",
			args: args{
				path: "/class/2/board",
				jwt:  jwtTest,
				params: &api.CreateBoardParams{
					Brdname: "mnewtest",
				},
				method: "POST",
			},
		},
		{
			name: "CreateComment",
			args: args{
				path: "/board/10_WhoAmI/article/" + string(articleID0) + "/comment",
				jwt:  jwtA1,
				params: &api.CreateCommentParams{
					CommentType: ptttype.COMMENT_TYPE_RECOMMEND,
					Content:     []byte("test123"),
				},
				method: "POST",
			},
		},
		{
			name: "GetArticle",
			args: args{
				path:   "/board/10_WhoAmI/article/" + string(articleID0),
				jwt:    jwtSysop,
				params: &api.GetArticleParams{},
				method: "GET",
			},
		},
		{
			name: "GetEmailTokenInfo",
			args: args{
				path:   "/emailtoken/info",
				jwt:    jwtSysop,
				params: &api.GetEmailTokenInfoParams{Jwt: jwtEmail, Context: api.CONTEXT_CHANGE_EMAIL},
				method: "POST",
			},
		},
		{
			name: "GetFavorites",
			args: args{
				path:   "/user/SYSOP/favorites",
				jwt:    jwtSysop,
				params: &api.GetFavoritesParams{},
				method: "GET",
			},
		},
		{
			name: "GetTokenInfo",
			args: args{
				path: "/token/info",
				jwt:  jwtSysop,
				params: &api.GetTokenInfoParams{
					Jwt: jwtSysop2,
				},
				method: "POST",
			},
		},
		{
			name: "GetUser",
			args: args{
				path:   "/user/SYSOP/information",
				jwt:    jwtSysop,
				params: &api.GetUserParams{},
				method: "GET",
			},
		},

		{
			name: "IsBoardValidUser",
			args: args{ // json: {}
				path:   "/board/10_WhoAmI/isvalid",
				jwt:    jwtSysop,
				params: nil,
				method: "GET",
			},
		},
		{
			name: "IsBoardsValidUser",
			args: args{ // json: {}
				path: "/boards/isvalid",
				jwt:  jwtSysop,
				params: &api.IsBoardsValidUserParams{
					BoardIDs: []bbs.BBoardID{"1_SYSOP"},
				},
				method: "POST",
			},
		},
		{
			name: "LoadAutoCompleteBoards",
			args: args{
				path: api.LOAD_AUTO_COMPLETE_BOARDS_R,
				jwt:  jwtSysop,
				params: &api.LoadAutoCompleteBoardsParams{
					StartIdx: "",
					NBoards:  4,
				},
				method: "GET",
			},
		},
		{
			name: "LoadBoardSummary",
			args: args{
				path:   "/board/6_ALLPOST/summary",
				jwt:    jwtSysop,
				params: &api.LoadBoardSummaryParams{},
				method: "GET",
			},
		},
		{
			name: "LoadBoardsByBids",
			args: args{
				path: "/board/6_ALLPOST/summary",
				jwt:  jwtSysop,
				params: &api.LoadBoardsByBidsParams{
					Bids: []ptttype.Bid{6, 7, 11, 8},
				},
				method: "GET",
			},
		},
		{
			name: "LoadBottomArticles",
			args: args{
				path:   "/board/10_WhoAmI/articles/bottom",
				jwt:    jwtSysop,
				method: "GET",
			},
		},
		{
			name: "LoadGeneralArticles",
			args: args{
				path: "/board/10_WhoAmI/articles",
				jwt:  jwtSysop,
				params: &api.LoadGeneralArticlesParams{
					StartIdx:  "",
					NArticles: 4,
				},
				method: "GET",
			},
		},
		{
			name: "LoadGeneralBoardsByClass",
			args: args{
				path: api.LOAD_GENERAL_BOARDS_BY_CLASS_R,
				jwt:  jwtSysop,
				params: &api.LoadGeneralBoardsParams{
					StartIdx: "",
					NBoards:  4,
					Asc:      true,
				},
				method: "GET",
			},
		},
		{
			name: "LoadGeneralBoards",
			args: args{
				path: api.LOAD_GENERAL_BOARDS_R,
				jwt:  jwtSysop,
				params: &api.LoadGeneralBoardsParams{
					StartIdx: "",
					NBoards:  4,
				},
				method: "GET",
			},
		},
		{
			name: "LoadHotBoards",
			args: args{
				path:   api.LOAD_HOT_BOARDS_R,
				jwt:    jwtSysop,
				method: "GET",
			},
		},
		{
			name: "Login",
			args: args{
				path: api.LOGIN_R,
				params: &api.LoginParams{
					Username: "SYSOP",
					Passwd:   "123123",
				},
				method: "POST",
			},
		},
		{
			name: "ReloadUHash",
			args: args{
				path:   "/admin/reloaduhash",
				jwt:    jwtSysop,
				method: "GET",
			},
		},
		{
			name: "SetIDEmail",
			args: args{
				path: "/user/SYSOP/setidemail",
				jwt:  jwtSysop,
				params: &api.SetIDEmailParams{
					Jwt:   jwtIDEmail,
					IsSet: true,
				},
				method: "POST",
			},
		},
		{
			name: "SetUserPerm",
			args: args{
				path: "/admin/user/SYSOP/setperm",
				jwt:  jwtSysop,
				params: &api.SetUserPermParams{
					Perm: ptttype.PERM_DEFAULT | ptttype.PERM_ADMIN | ptttype.PERM_LOGINOK,
				},
				method: "POST",
			},
		},
		{
			name: "SetUserPerm-2",
			args: args{
				path: "/admin/user/SYSOP/setperm",
				jwt:  jwtSysop,
				params: &api.SetUserPermParams{
					Perm: testUserec.Userlevel,
				},
				method: "POST",
			},
		},
		{
			name: "CrossPost",
			args: args{
				path: fmt.Sprintf("/board/%v/article/%v/crosspost", boardSummary.BBoardID, articleSummary0.ArticleID),
				jwt:  jwtSysop,
				params: &api.CrossPostParams{
					XBoardID: forwardBoardSummary.BBoardID,
				},
				method: "POST",
			},
		},
		{
			name: "LoadClassBoards",
			args: args{
				path: fmt.Sprintf("/cls/1/boards"),
				jwt:  jwtSysop,
				params: &api.LoadClassBoardsParams{
					IsSystem: true,
				},
				method: "GET",
			},
		},
		{
			name: "LoadFullClassBoards",
			args: args{
				path: fmt.Sprintf("/cls/boards"),
				jwt:  jwtSysop,
				params: &api.LoadFullClassBoardsParams{
					StartBid: 1,
					NBoards:  100,
					IsSystem: true,
				},
				method: "GET",
			},
		},
		{
			name: "LoadBoardDetail",
			args: args{
				path:   fmt.Sprintf("/board/10_WhoAmI/detail"),
				jwt:    jwtSysop,
				method: "GET",
			},
		},
		{
			name: "GetUserVisitCount",
			args: args{
				path:   fmt.Sprintf("/uservisitcount"),
				jwt:    jwtTest,
				method: "GET",
			},
		},
		{
			name: "EditArticle",
			args: args{
				path: fmt.Sprintf("/board/%v/article/%v/edit", boardSummary.BBoardID, articleSummary0.ArticleID),
				jwt:  jwtTest,
				params: &api.EditArticleParams{
					Content: editContent0,
					OldSZ:   len(gotContent0),
					OldSum:  oldSum0,
				},
				method: "POST",
			},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			w := httptest.NewRecorder()
			req := setRequest(tt.args.path, tt.args.params, tt.args.jwt, nil, tt.args.method)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}
		})
		wg.Wait()
	}
}
