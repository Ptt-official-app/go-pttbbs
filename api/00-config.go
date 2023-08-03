package api

import (
	"time"
)

var (
	// Creating JWT Token
	JWT_ISSUER = "go-pttbbs"
	GUEST      = "guest"

	JWT_SECRET                = []byte("jwt_secret")
	JWT_TOKEN_EXPIRE_TS       = 86400 * 1 // 1 days
	JWT_TOKEN_EXPIRE_DURATION = time.Duration(JWT_TOKEN_EXPIRE_TS) * time.Second

	EMAIL_JWT_SECRET                = []byte("email_jwt_secret")
	EMAIL_JWT_TOKEN_EXPIRE_TS       = 60 * 15 // 15 mins
	EMAIL_JWT_TOKEN_EXPIRE_DURATION = time.Duration(EMAIL_JWT_TOKEN_EXPIRE_TS) * time.Second

	REFRESH_JWT_CLAIM_TYPE            = "refresh"
	REFRESH_JWT_SECRET                = []byte("refresh_jwt_secret")
	REFRESH_JWT_TOKEN_EXPIRE_TS       = 86400 * 7 // 7 days
	REFRESH_JWT_TOKEN_EXPIRE_DURATION = time.Duration(REFRESH_JWT_TOKEN_EXPIRE_TS) * time.Second
)
