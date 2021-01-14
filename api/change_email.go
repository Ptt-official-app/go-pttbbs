package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const CHANGE_EMAIL_R = "/user/:uid/changeemail"

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

//ChangeEmail
//
//Sysop initiates only attempt-change-mail.
//Sysop does not change email directly.
func ChangeEmail(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	theParams, ok := params.(*ChangeEmailParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*ChangeEmailPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	isValid, email := userInfoIsValidEmailUser(uuserID, thePath.UserID, theParams.Jwt, CONTEXT_CHANGE_EMAIL, false)
	if !isValid {
		return nil, ErrInvalidUser
	}

	err = bbs.ChangeEmail(thePath.UserID, email)
	if err != nil {
		return nil, err
	}

	result = &ChangeEmailResult{
		UserID: thePath.UserID,
		Email:  email,
	}

	return result, nil
}
