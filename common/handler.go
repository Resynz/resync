// Package common /**
package common

import (
	"net/http"
	"resync/code"
)

func HandleResponse(ctx *Context, c code.ResponseCode, d interface{}, msg ...string) {
	m := code.GetCodeMsg(c)

	if len(msg) > 0 {
		m = msg[0]
	}
	data := map[string]interface{}{
		"code":       c,
		"message":    m,
	}
	if d != nil {
		data["data"] = d
	}
	ctx.JSON(http.StatusOK, data)
}
