package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

func processResult(c *gin.Context, result interface{}, err error) {
	switch err {
	case nil:
		c.JSON(200, result)

	case ErrInvalidHost:
		c.JSON(400, &errResult{err.Error()})
	case ErrInvalidRemoteAddr:
		c.JSON(400, &errResult{err.Error()})

	case ErrInvalidParams:
		c.JSON(400, &errResult{err.Error()})
	case ErrInvalidPath:
		c.JSON(400, &errResult{err.Error()})

	case bbs.ErrInvalidParams:
		c.JSON(400, &errResult{err.Error()})

	case ptttype.ErrUserIDAlreadyExists:
		c.JSON(400, &errResult{err.Error()})

	//401
	case ErrInvalidToken:
		c.JSON(401, &errResult{err.Error()})

	case ErrInvalidToken:
		c.JSON(401, &errResult{err.Error()})
	case ErrLoginFailed:
		c.JSON(401, &errResult{err.Error()})

	//403
	case ErrInvalidUser:
		c.JSON(403, &errResult{err.Error()})
	case ptttype.ErrInvalidUserID:
		c.JSON(403, &errResult{err.Error()})

	default:
		c.JSON(500, &errResult{err.Error()})
	}
}
