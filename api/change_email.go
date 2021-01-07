package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const CHANGE_EMAIL_R = "/users/:uid/changeemail"

type ChangeEmailParams struct {
	Jwt string `json:"email_token"`
}

type ChangeEmailPath struct {
	UserID bbs.UUserID `uri:"uid" binding:"required"`
}

type ChangeEmailResult struct {
	UserID bbs.UUserID `json:"user_id"`
	Email  string      `json:"email"`
}

func ChangeEmailWrapper(c *gin.Context) {
	params := &ChangeEmailParams{}
	path := &ChangeEmailPath{}
	LoginRequiredPathJSON(ChangeEmail, params, path, c)
}

func ChangeEmail(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	theParams, ok := params.(*ChangeEmailParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*ChangeEmailPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	if uuserID != thePath.UserID {
		return nil, ErrInvalidUser
	}

	userID, _, email, err := VerifyEmailJwt(theParams.Jwt)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if uuserID != userID {
		return nil, ErrInvalidUser
	}

	err = bbs.ChangeEmail(userID, email)
	if err != nil {
		return nil, err
	}

	result = &ChangeEmailResult{
		UserID: userID,
		Email:  email,
	}

	return result, nil
}
