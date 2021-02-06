package main

import (
	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	apiPrefix = "/v1"
)

func withPrefix(path string) string {
	return apiPrefix + path
}

func initGin() (*gin.Engine, error) {
	router := gin.Default()

	//options
	router.OPTIONS("/*path", api.OptionsWrapper)

	router.POST(withPrefix(api.INDEX_R), api.IndexWrapper)

	//login/register
	router.POST(withPrefix(api.LOGIN_R), api.LoginWrapper)
	router.POST(withPrefix(api.REGISTER_R), api.RegisterWrapper)

	//board
	router.GET(withPrefix(api.LOAD_GENERAL_BOARDS_R), api.LoadGeneralBoardsWrapper)
	router.GET(withPrefix(api.LOAD_GENERAL_ARTICLES_R), api.LoadGeneralArticlesWrapper)
	router.GET(withPrefix(api.LOAD_BOARD_SUMMARY_R), api.LoadBoardSummaryWrapper)
	router.GET(withPrefix(api.LOAD_HOT_BOARDS_R), api.LoadHotBoardsWrapper)
	router.GET(withPrefix(api.LOAD_GENERAL_BOARDS_BY_CLASS_R), api.LoadGeneralBoardsByClassWrapper)
	router.POST(withPrefix(api.LOAD_BOARDS_BY_BIDS_R), api.LoadBoardsByBidsWrapper)
	router.GET(withPrefix(api.IS_BOARD_VALID_USER_R), api.IsBoardValidUserWrapper)
	router.GET(withPrefix(api.LOAD_AUTO_COMPLETE_BOARDS_R), api.LoadAutoCompleteBoardsWrapper)

	//article
	router.GET(withPrefix(api.GET_ARTICLE_R), api.GetArticleWrapper)

	//user
	router.GET(withPrefix(api.GET_USER_R), api.GetUserWrapper)
	router.POST(withPrefix(api.CHANGE_PASSWD_R), api.ChangePasswdWrapper)
	router.POST(withPrefix(api.CHANGE_EMAIL_R), api.ChangeEmailWrapper)
	router.POST(withPrefix(api.ATTEMPT_CHANGE_EMAIL_R), api.AttemptChangeEmailWrapper)
	router.POST(withPrefix(api.ATTEMPT_SET_ID_EMAIL_R), api.AttemptSetIDEmailWrapper)
	router.POST(withPrefix(api.SET_ID_EMAIL_R), api.SetIDEmailWrapper)
	router.POST(withPrefix(api.GET_TOKEN_INFO_R), api.GetTokenInfoWrapper)
	router.POST(withPrefix(api.GET_EMAIL_TOKEN_INFO_R), api.GetEmailTokenInfoWrapper)
	router.GET(withPrefix(api.GET_FAV_R), api.GetFavoritesWrapper)
	router.POST(withPrefix(api.CHECK_EXISTS_USER_R), api.CheckExistsUserWrapper)

	return router, nil
}

func main() {
	err := initMain()
	if err != nil {
		log.Errorf("unable to initMain: e: %v", err)
		return
	}
	router, err := initGin()
	if err != nil {
		return
	}

	_ = router.Run(HTTP_HOST)
}
