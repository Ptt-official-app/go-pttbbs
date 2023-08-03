package api

import (
	"github.com/gin-gonic/gin"
)

const REFRESH_R = "/refresh"

type RefreshParams struct {
	ClientInfo string `json:"client_info"`
	Refresh    string `json:"refresh_token"`
}

type RefreshResult LoginResult

func RefreshWrapper(c *gin.Context) {
	params := &RefreshParams{}
	JSON(Refresh, params, c)
}

func Refresh(remoteAddr string, params interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*RefreshParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	jwt := GetJwt(c)

	jwtUserID, jwtExpireTS, jwtClientInfo, err := VerifyJwt(jwt, false)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, refreshExpireTS, clientInfo, err := VerifyRefreshJwt(theParams.Refresh)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// verify that jwt and refresh-jwt are with same pair.
	diffExpireTS := refreshExpireTS - jwtExpireTS
	expectedDiffExpireTS := REFRESH_JWT_TOKEN_EXPIRE_TS - JWT_TOKEN_EXPIRE_TS
	diffDiffExpireTS := diffExpireTS - expectedDiffExpireTS
	if diffDiffExpireTS > EPSILON_EXPIRE_TS || diffDiffExpireTS < -EPSILON_EXPIRE_TS {
		return nil, ErrInvalidToken
	}

	if clientInfo != theParams.ClientInfo && clientInfo != jwtClientInfo {
		return nil, ErrInvalidToken
	}

	if userID != jwtUserID {
		return nil, ErrInvalidToken
	}

	token, err := CreateToken(userID, clientInfo)
	if err != nil {
		return nil, err
	}

	refreshToken, err := CreateRefreshToken(userID, clientInfo)
	if err != nil {
		return nil, err
	}

	result = &RefreshResult{
		UserID:    userID,
		Jwt:       token,
		TokenType: "bearer",
		Refresh:   refreshToken,
	}

	return result, nil
}
