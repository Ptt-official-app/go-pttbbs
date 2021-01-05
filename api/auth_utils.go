package api

import (
	"strings"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func GetJwt(c *gin.Context) (jwt string) {
	tokenStr := strings.TrimSpace(c.GetHeader("Authorization"))
	tokenList := strings.Split(tokenStr, " ")
	if len(tokenList) != 2 {
		return ""
	}

	return tokenList[1]
}

func VerifyJwt(raw string) (userID bbs.UUserID, clientInfo string, err error) {
	if raw == "" {
		return bbs.UUserID(GUEST), "", nil
	}

	tok, err := jwt.ParseSigned(raw)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	cl := &JwtClaim{}
	if err := tok.Claims(JWT_SECRET, cl); err != nil {
		return "", "", ErrInvalidToken
	}

	if cl.Expire == nil {
		return "", "", ErrInvalidToken
	}

	currentNanoTS := jwt.NewNumericDate(time.Now())
	if *currentNanoTS > *cl.Expire {
		return "", "", ErrInvalidToken
	}

	return cl.UUserID, cl.ClientInfo, nil
}

func createToken(userID bbs.UUserID, clientInfo string) (string, error) {
	var err error

	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: JWT_SECRET}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}

	cl := &JwtClaim{
		ClientInfo: clientInfo,
		UUserID:    userID,
		Expire:     jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 2)),
	}

	raw, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		return "", err
	}

	return raw, nil
}
