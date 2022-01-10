/**
 * @Author: Resynz
 * @Date: 2022/1/5 13:57
 */
package admin

import (
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
	"time"
)

// Login 登录
func Login(ctx *common.Context) {
	type formValidate struct {
		Name string `form:"name" binding:"required" json:"name"`
		Password string `form:"password" binding:"required" json:"password"`
	}
	var form formValidate
	if err:=ctx.ShouldBind(&form);err!=nil{
		common.HandleResponse(ctx,code.InvalidParams,nil)
		return
	}
	var admin model.Admin
	has,err:=db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("name",form.Name),dbx.Eq("password",form.Password)).Get(&admin)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx,code.BadRequest,nil,"账号或密码错误")
		return
	}
	auth:=&common.Auth{
		Id:     admin.Id,
		Name:   admin.Name,
		Ip:     ctx.ClientIP(),
		Expire: 0,
	}
	token,err:=auth.GenToken()
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	data:=map[string]string{
		"token":token,
	}

	loginLog:=model.LoginLog{
		Id:         0,
		AdminId:    admin.Id,
		Ip:         auth.Ip,
		UserAgent:  ctx.Request.UserAgent(),
		CreateTime: time.Now().Unix(),
	}

	_ = db.Handler.XStmt(loginLog.GetTableName()).Insert(&loginLog)

	common.HandleResponse(ctx,code.SuccessCode,data)
}
