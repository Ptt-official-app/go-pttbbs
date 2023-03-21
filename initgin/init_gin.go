package initgin

import (
	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/gin-gonic/gin"
)

var apiPrefix = "/v1"

func withPrefix(path string) string {
	return apiPrefix + path
}

func InitGin() (*gin.Engine, error) {
	router := gin.Default()

	// options
	router.OPTIONS("/*path", api.OptionsWrapper)

	router.POST(withPrefix(api.INDEX_R), api.IndexWrapper)
	router.GET(withPrefix(api.GET_VERSION_R), api.GetVersionWrapper)

	// login/register
	router.POST(withPrefix(api.LOGIN_R), api.LoginWrapper)
	router.POST(withPrefix(api.REGISTER_R), api.RegisterWrapper)

	// board
	router.GET(withPrefix(api.LOAD_GENERAL_BOARDS_R), api.LoadGeneralBoardsWrapper)
	router.GET(withPrefix(api.LOAD_GENERAL_ARTICLES_R), api.LoadGeneralArticlesWrapper)
	router.GET(withPrefix(api.LOAD_BOTTOM_ARTICLES_R), api.LoadBottomArticlesWrapper)
	router.GET(withPrefix(api.LOAD_BOARD_SUMMARY_R), api.LoadBoardSummaryWrapper)
	router.GET(withPrefix(api.LOAD_HOT_BOARDS_R), api.LoadHotBoardsWrapper)
	router.GET(withPrefix(api.LOAD_GENERAL_BOARDS_BY_CLASS_R), api.LoadGeneralBoardsByClassWrapper)
	router.POST(withPrefix(api.LOAD_BOARDS_BY_BIDS_R), api.LoadBoardsByBidsWrapper)
	router.GET(withPrefix(api.IS_BOARD_VALID_USER_R), api.IsBoardValidUserWrapper)
	router.POST(withPrefix(api.IS_BOARDS_VALID_USER_R), api.IsBoardsValidUserWrapper)
	router.GET(withPrefix(api.LOAD_AUTO_COMPLETE_BOARDS_R), api.LoadAutoCompleteBoardsWrapper)
	router.GET(withPrefix(api.LOAD_BOARD_DETAIL_R), api.LoadBoardDetailWrapper)

	router.POST(withPrefix(api.CREATE_BOARD_R), api.CreateBoardWrapper)

	router.GET(withPrefix(api.LOAD_CLASS_BOARDS_R), api.LoadClassBoardsWrapper)
	router.GET(withPrefix(api.LOAD_FULL_CLASS_BOARDS_R), api.LoadFullClassBoardsWrapper)

	router.GET(withPrefix(api.LOAD_GENERAL_BOARD_DETAILS_R), api.LoadGeneralBoardDetailsWrapper)

	router.GET(withPrefix(api.GET_POST_TEMPLATE_R), api.GetPostTemplateWrapper)

	// article
	router.GET(withPrefix(api.GET_ARTICLE_R), api.GetArticleWrapper)
	router.POST(withPrefix(api.CREATE_ARTICLE_R), api.CreateArticleWrapper)
	router.POST(withPrefix(api.CROSS_POST_R), api.CrossPostWrapper)
	router.POST(withPrefix(api.EDIT_ARTICLE_R), api.EditArticleWrapper)
	router.POST(withPrefix(api.DELETE_ARTICLES_R), api.DeleteArticlesWrapper)

	// comment
	router.POST(withPrefix(api.CREATE_COMMENT_R), api.CreateCommentWrapper)

	// user
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
	router.GET(withPrefix(api.GET_USER_VISIT_COUNT_R), api.GetUserVisitCountWrapper)

	// admin
	router.GET(withPrefix(api.RELOAD_UHASH_R), api.ReloadUHashWrapper)
	router.POST(withPrefix(api.SET_USER_PERM_R), api.SetUserPermWrapper)

	return router, nil
}
