package api

import (
	"time"

	configutil "github.com/Ptt-official-app/go-pttbbs/configutil"
)

const configPrefix = "go-pttbbs:api"

func InitConfig() error {
	config()
	postInitConfig()
	return nil
}

func setStringConfig(idx string, orig string) string {
	return configutil.SetStringConfig(configPrefix, idx, orig)
}

func setBytesConfig(idx string, orig []byte) []byte {
	return configutil.SetBytesConfig(configPrefix, idx, orig)
}

func setIntConfig(idx string, orig int) int {
	return configutil.SetIntConfig(configPrefix, idx, orig)
}

func postInitConfig() {
	_ = setJwtTokenExpireTS(JWT_TOKEN_EXPIRE_TS)
	_ = setEmailJwtTokenExpireTS(EMAIL_JWT_TOKEN_EXPIRE_TS)
	_ = setRefreshJwtTokenExpireTS(REFRESH_JWT_TOKEN_EXPIRE_TS)
}

func setJwtTokenExpireTS(jwtTokenExpireTS int) (origJwtTokenExpireTS int) {
	origJwtTokenExpireTS = JWT_TOKEN_EXPIRE_TS

	JWT_TOKEN_EXPIRE_TS = jwtTokenExpireTS
	JWT_TOKEN_EXPIRE_DURATION = time.Duration(JWT_TOKEN_EXPIRE_TS) * time.Second

	return origJwtTokenExpireTS
}

func setEmailJwtTokenExpireTS(emailJwtTokenExpireTS int) (origEmailJwtTokenExpireTS int) {
	origEmailJwtTokenExpireTS = EMAIL_JWT_TOKEN_EXPIRE_TS

	EMAIL_JWT_TOKEN_EXPIRE_TS = emailJwtTokenExpireTS
	EMAIL_JWT_TOKEN_EXPIRE_DURATION = time.Duration(EMAIL_JWT_TOKEN_EXPIRE_TS) * time.Second

	return origEmailJwtTokenExpireTS
}

func setRefreshJwtTokenExpireTS(refreshJwtTokenExpireTS int) (origRefreshJwtTokenExpireTS int) {
	origRefreshJwtTokenExpireTS = REFRESH_JWT_TOKEN_EXPIRE_TS

	REFRESH_JWT_TOKEN_EXPIRE_TS = refreshJwtTokenExpireTS
	REFRESH_JWT_TOKEN_EXPIRE_DURATION = time.Duration(REFRESH_JWT_TOKEN_EXPIRE_TS) * time.Second

	return origRefreshJwtTokenExpireTS
}
