package api

import "gopkg.in/square/go-jose.v2/jwt"

type ApiFunc func(remoteAddr string, params interface{}) (interface{}, error)

type LoginRequiredApiFunc func(remoteAddr string, userID string, params interface{}) (interface{}, error)

type LoginRequiredPathApiFunc func(remoteAddr string, userID string, params interface{}, path interface{}) (interface{}, error)

type JwtClaim struct {
	UserID string
	Expire *jwt.NumericDate
}
