package api

import (
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/types"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GetJwt(c *gin.Context) (jwt string) {
	tokenStr := strings.TrimSpace(c.GetHeader("Authorization"))
	tokenList := strings.Split(tokenStr, " ")
	if len(tokenList) != 2 {
		return ""
	}

	return tokenList[1]
}

func VerifyJwt(raw string, isCheckExpire bool) (userID bbs.UUserID, expireTS int, clientInfo string, err error) {
	if raw == "" {
		return bbs.UUserID(GUEST), 0, "", nil
	}

	cl, err := parseJwtClaim(raw)
	if err != nil {
		return "", 0, "", ErrInvalidToken
	}

	if isCheckExpire {
		currentTS := int(types.NowTS())
		if currentTS > cl.Expire {
			return "", 0, "", ErrInvalidToken
		}
	}

	return bbs.UUserID(cl.UUserID), cl.Expire, cl.ClientInfo, nil
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

func CreateToken(userID bbs.UUserID, clientInfo string) (raw string, expireTime types.Time4, err error) {
	defer func() {
		err2 := recover()
		if err2 == nil {
			return
		}

		err = types.ErrRecover(err2)
	}()

	expireTime = types.NowTS() + types.Time4(JWT_TOKEN_EXPIRE_TS)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"cli": clientInfo,
		"sub": userID,
		"exp": int(expireTime),
	})

	raw, err = token.SignedString(JWT_SECRET)
	if err != nil {
		return "", 0, err
	}

	return raw, expireTime, nil
}

func VerifyEmailJwt(raw string, context EmailTokenContext) (userID bbs.UUserID, expireTS int, clientInfo string, email string, err error) {
	if raw == "" {
		return "", 0, "", "", ErrInvalidToken
	}

	cl, err := parseEmailJwtClaim(raw)
	if err != nil {
		return "", 0, "", "", ErrInvalidToken
	}

	currentTS := int(types.NowTS())
	if currentTS > cl.Expire {
		return "", 0, "", "", ErrInvalidToken
	}

	if cl.Context != string(context) {
		return "", 0, "", "", ErrInvalidToken
	}

	return bbs.UUserID(cl.UUserID), cl.Expire, cl.ClientInfo, cl.Email, nil
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

func VerifyRefreshJwt(raw string) (userID bbs.UUserID, expireTS int, clientInfo string, err error) {
	if raw == "" {
		return bbs.UUserID(GUEST), 0, "", nil
	}

	cl, err := parseRefreshJwtClaim(raw)
	if err != nil {
		return "", 0, "", ErrInvalidToken
	}

	currentTS := int(types.NowTS())
	if currentTS > cl.Expire {
		return "", 0, "", ErrInvalidToken
	}

	if cl.TheType != REFRESH_JWT_CLAIM_TYPE {
		return "", 0, "", ErrInvalidToken
	}

	return bbs.UUserID(cl.UUserID), cl.Expire, cl.ClientInfo, nil
}

func parseRefreshJwtClaim(raw string) (cl *RefreshJwtClaim, err error) {
	tok, err := ParseJwt(raw, REFRESH_JWT_SECRET)
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
	typ, err := ParseClaimString(claim, "typ")
	if err != nil {
		return nil, err
	}

	cl = &RefreshJwtClaim{
		ClientInfo: cli,
		UUserID:    sub,
		Expire:     exp,
		TheType:    typ,
	}

	return cl, nil
}

func CreateRefreshToken(userID bbs.UUserID, clientInfo string) (raw string, expireTime types.Time4, err error) {
	defer func() {
		err2 := recover()
		if err2 == nil {
			return
		}

		err = types.ErrRecover(err2)
	}()

	expireTime = types.NowTS() + types.Time4(REFRESH_JWT_TOKEN_EXPIRE_TS)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"cli": clientInfo,
		"sub": userID,
		"exp": int(expireTime),
		"typ": REFRESH_JWT_CLAIM_TYPE,
	})

	raw, err = token.SignedString(REFRESH_JWT_SECRET)
	if err != nil {
		return "", 0, err
	}

	return raw, expireTime, nil
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
