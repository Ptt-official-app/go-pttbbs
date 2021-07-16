package api

import (
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

	cl, err := parseJwtClaim(raw)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	currentTS := int(types.NowTS())
	if currentTS > cl.Expire {
		return "", "", ErrInvalidToken
	}

	return bbs.UUserID(cl.UUserID), cl.ClientInfo, nil
}

func parseJwtClaim(raw string) (cl *JwtClaim, err error) {
	tok, err := ParseJwt(raw, JWT_SECRET)
	if err != nil {
		return nil, err
	}

	claim, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	cli, err := ParseClaimString(claim, "cli")
	if err != nil {
		return nil, err
	}
	sub, err := ParseClaimString(claim, "sub")
	if err != nil {
		return nil, err
	}
	exp, err := ParseClaimInt(claim, "exp")
	if err != nil {
		return nil, err
	}

	cl = &JwtClaim{
		ClientInfo: cli,
		UUserID:    sub,
		Expire:     exp,
	}

	return cl, nil
}

func CreateToken(userID bbs.UUserID, clientInfo string) (raw string, err error) {
	defer func() {
		err2 := recover()
		if err2 == nil {
			return
		}

		err = types.ErrRecover(err2)
	}()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"cli": clientInfo,
		"sub": userID,
		"exp": int(types.NowTS()) + JWT_TOKEN_EXPIRE_TS,
	})

	raw, err = token.SignedString(JWT_SECRET)
	if err != nil {
		return "", err
	}

	return raw, nil
}

func VerifyEmailJwt(raw string, context EmailTokenContext) (userID bbs.UUserID, clientInfo string, email string, err error) {
	if raw == "" {
		return "", "", "", ErrInvalidToken
	}

	cl, err := parseEmailJwtClaim(raw)
	if err != nil {
		return "", "", "", ErrInvalidToken
	}

	currentTS := int(types.NowTS())
	if currentTS > cl.Expire {
		return "", "", "", ErrInvalidToken
	}

	if cl.Context != string(context) {
		return "", "", "", ErrInvalidToken
	}

	return bbs.UUserID(cl.UUserID), cl.ClientInfo, cl.Email, nil
}

func parseEmailJwtClaim(raw string) (cl *EmailJwtClaim, err error) {
	tok, err := ParseJwt(raw, EMAIL_JWT_SECRET)
	if err != nil {
		return nil, err
	}

	claim, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	cli, err := ParseClaimString(claim, "cli")
	if err != nil {
		return nil, err
	}
	sub, err := ParseClaimString(claim, "sub")
	if err != nil {
		return nil, err
	}
	eml, err := ParseClaimString(claim, "eml")
	if err != nil {
		return nil, err
	}
	exp, err := ParseClaimInt(claim, "exp")
	if err != nil {
		return nil, err
	}
	ctx, err := ParseClaimString(claim, "ctx")
	if err != nil {
		return nil, err
	}

	cl = &EmailJwtClaim{
		ClientInfo: cli,
		UUserID:    sub,
		Email:      eml,
		Expire:     exp,
		Context:    ctx,
	}

	return cl, nil
}

func CreateEmailToken(userID bbs.UUserID, clientInfo string, email string, context EmailTokenContext) (raw string, err error) {
	defer func() {
		err2 := recover()
		if err2 == nil {
			return
		}

		err = types.ErrRecover(err2)
	}()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"cli": clientInfo,
		"sub": userID,
		"eml": email,
		"exp": int(types.NowTS()) + JWT_TOKEN_EXPIRE_TS,
		"ctx": context,
	})

	raw, err = token.SignedString(EMAIL_JWT_SECRET)
	if err != nil {
		return "", err
	}

	return raw, nil
}

func ParseJwt(raw string, secret []byte) (tok *jwt.Token, err error) {
	tok, err = jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	return tok, err
}

func ParseClaimString(claim jwt.MapClaims, idx string) (ret string, err error) {
	ret_i, ok := claim[idx]
	if !ok {
		return "", nil
	}
	ret, ok = ret_i.(string)
	if !ok {
		return "", ErrInvalidToken
	}

	return ret, nil
}

func ParseClaimInt(claim jwt.MapClaims, idx string) (ret int, err error) {
	ret_i, ok := claim[idx]
	if !ok {
		return 0, nil
	}
	// XXX it's float64 in go-jwt, but it's ok to have second(time)-level inaccuracy for expire-ts.
	ret_f64, ok := ret_i.(float64)
	if !ok {
		return 0, ErrInvalidToken
	}

	return int(ret_f64), nil
}
