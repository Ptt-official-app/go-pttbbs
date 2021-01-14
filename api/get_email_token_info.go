package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const GET_EMAIL_TOKEN_INFO_R = "/emailtoken/info"

type GetEmailTokenInfoParams struct {
	Jwt     string            `json:"token" form:"token" url:"token"`
	Context EmailTokenContext `json:"context" form:"context" url:"context"`
}

type GetEmailTokenInfoResult struct {
	ClientInfo string      `json:"client_info"`
	UserID     bbs.UUserID `json:"user_id"`
	Email      string      `json:"email"`
}

func GetEmailTokenInfoWrapper(c *gin.Context) {
	params := &GetEmailTokenInfoParams{}

	LoginRequiredJSON(GetEmailTokenInfo, params, c)
}

func GetEmailTokenInfo(remoteAddr string, uuserID bbs.UUserID, params interface{}) (result interface{}, err error) {
	theParams, ok := params.(*GetEmailTokenInfoParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	userID, clientInfo, email, err := VerifyEmailJwt(theParams.Jwt, theParams.Context)
	if err != nil {
		return nil, err
	}

	isValid, _ := userInfoIsValidEmailUser(uuserID, userID, theParams.Jwt, theParams.Context, true)
	if !isValid {
		return nil, ErrInvalidUser
	}

	result = &GetEmailTokenInfoResult{
		ClientInfo: clientInfo,
		UserID:     userID,
		Email:      email,
	}

	return result, nil
}
