package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const GET_REFRESH_TOKEN_INFO_R = "/refreshtoken/info"

type GetRefreshTokenInfoParams struct {
	Jwt string `json:"token" form:"token" url:"token"`
}

type GetRefreshTokenInfoResult struct {
	ClientInfo string      `json:"client_info"`
	UserID     bbs.UUserID `json:"user_id"`
	Expire     int         `json:"expire"`
}

func GetRefreshTokenInfoWrapper(c *gin.Context) {
	params := &GetRefreshTokenInfoParams{}

	LoginRequiredJSON(GetRefreshTokenInfo, params, c)
}

func GetRefreshTokenInfo(remoteAddr string, uuserID bbs.UUserID, params interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*GetRefreshTokenInfoParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	userID, expireTS, clientInfo, err := VerifyRefreshJwt(theParams.Jwt)
	if err != nil {
		return nil, err
	}

	if userID != uuserID {
		return nil, ErrInvalidToken
	}

	result = &GetRefreshTokenInfoResult{
		ClientInfo: clientInfo,
		UserID:     userID,
		Expire:     expireTS,
	}

	return result, nil
}
