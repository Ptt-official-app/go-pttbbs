package api

import "time"

var (
	//Creating JWT Token
	JWT_SECRET = []byte("jwt_secret")
	JWT_ISSUER = "go-pttbbs"
	GUEST      = "guest"

	EMAIL_JWT_SECRET = []byte("email_jwt_secret")

	JWT_TOKEN_EXPIRE_TS             = 3600 * 24 * 1 // 1 days
	JWT_TOKEN_EXPIRE_DURATION       = time.Duration(JWT_TOKEN_EXPIRE_TS) * time.Second
	EMAIL_JWT_TOKEN_EXPIRE_TS       = 60 * 15 // 15 mins
	EMAIL_JWT_TOKEN_EXPIRE_DURATION = time.Duration(EMAIL_JWT_TOKEN_EXPIRE_TS) * time.Second
)
