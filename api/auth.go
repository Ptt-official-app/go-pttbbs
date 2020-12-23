package api

import (
	"time"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"gopkg.in/square/go-jose.v2/jwt"
)

func VerifyJwt(raw string) (userID bbs.UUserID, err error) {
	if raw == "" {
		return bbs.UUserID(GUEST), nil
	}

	tok, err := jwt.ParseSigned(raw)
	if err != nil {
		return "", ErrInvalidToken
	}

	cl := &JwtClaim{}
	if err := tok.Claims(JWT_SECRET, cl); err != nil {
		return "", ErrInvalidToken
	}

	currentNanoTS := jwt.NewNumericDate(time.Now())
	if *currentNanoTS > *cl.Expire {
		return "", ErrInvalidToken
	}

	return cl.UUserID, nil
}
