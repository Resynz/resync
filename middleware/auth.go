// Package middleware /**
package middleware

import (
	"github.com/gin-gonic/gin"
	"resync/common"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ak string
		if a := c.GetHeader("token"); a != "" {
			ak = a
		}
		if ak == "" {
			if a := c.Query("token"); a != "" {
				ak = a
			}
		}
		if ak != "" {
			auth := &common.Auth{}
			if err := auth.ParseToken(ak); err == nil {
				c.Set("Auth-Id", auth.Id)
				c.Set("Auth-Ip", auth.Ip)
				c.Set("Auth-Expire", auth.Expire)
				c.Set("Auth-Name",auth.Name)
			}
		}
	}
}
