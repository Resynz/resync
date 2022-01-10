// Package common /**
package common

import (
	"github.com/gin-gonic/gin"
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/db"
	"resync/db/model"
	"resync/enum"
	"time"
)

type Context struct {
	*gin.Context
	Auth *Auth
}

type HandlerFunc func(ctx *Context)

func AuthDetection(next HandlerFunc, rules ...enum.PermissionEnum) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		context := new(Context)
		context.Context = ctx
		context.Auth = new(Auth)
		authId := ctx.GetInt64("Auth-Id")
		if authId > 0 {
			authIp := ctx.GetString("Auth-Ip")
			authExpire := ctx.GetInt64("Auth-Expire")
			if authExpire <= time.Now().Unix() {
				HandleResponse(context, code.InvalidAuth, nil, "登录超时")
				return
			}
			if authIp != ctx.ClientIP() {
				HandleResponse(context, code.InvalidAuth, nil, "设备IP变更，请重新登录")
				return
			}
			var admin model.Admin
			has,err:=db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id",authId)).Get(&admin)
			if err!=nil{
				HandleResponse(context,code.BadRequest,nil,err.Error())
				return
			}
			if !has {
				HandleResponse(context, code.InvalidAuth, nil, "账号不存在")
				return
			}
			if admin.Status != model.AccountStatusEnable {
				HandleResponse(context,code.InvalidAuth,nil,"账号已禁用")
				return
			}
			context.Auth.Id = authId
			context.Auth.Ip = authIp
			context.Auth.Expire = authExpire
			context.Auth.Name = ctx.GetString("Auth-Name")
			if authExpire-time.Now().Unix() <= 30*60 {
				token, err := context.Auth.GenToken()
				if err != nil {
					HandleResponse(context, code.BadRequest, nil, err.Error())
					return
				}
				context.Header("X-Refresh-Token", token)
			}
		}
		if len(rules) != 0 {
			for _, r := range rules {
				if r == enum.MustLogin {
					if context.Auth.Id == 0 {
						HandleResponse(context, code.InvalidAuth, nil)
						return
					}
				}
			}
		}
		next(context)
	}
}
