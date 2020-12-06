package api

import "github.com/Ptt-official-app/go-pttbbs/bbs"

type RegisterParams struct {
	UserID string
	Passwd string
	IP     string
	Email  string

	Nickname string
	Realname string
	Career   string
	Address  string
	Over18   bool
}

type RegisterResult struct {
	Jwt string
}

func Register(params interface{}) (interface{}, error) {
	registerParams, ok := params.(*RegisterParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	user, err := bbs.Register(
		registerParams.UserID,
		registerParams.Passwd,
		registerParams.IP,
		registerParams.Email,

		registerParams.Nickname,
		registerParams.Realname,
		registerParams.Career,
		registerParams.Address,
		registerParams.Over18,
	)
	if err != nil {
		return nil, err
	}

	token, err := createToken(user)
	if err != nil {
		return nil, err
	}

	result := &RegisterResult{
		Jwt: token,
	}

	return result, nil
}
