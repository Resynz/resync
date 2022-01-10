/**
 * @Author: Resynz
 * @Date: 2022/1/7 11:01
 */
package code_auth

import (
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
)

func Update(ctx *common.Context) {
	var codeAuth model.CodeAuth
	has,err:=db.Handler.XStmt(codeAuth.GetTableName()).Where(dbx.Eq("id",ctx.Param("id"))).Get(&codeAuth)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx,code.InvalidRequest,nil,"该账号不存在")
		return
	}
	type formValidate struct {
		UserName string `form:"user_name" binding:"required" json:"user_name"`
		Password string `form:"password" binding:"" json:"password"`
	}
	var form formValidate
	if err = ctx.ShouldBind(&form);err!=nil{
		common.HandleResponse(ctx,code.InvalidParams,nil)
		return
	}
	codeAuth.UserName = form.UserName
	if form.Password != "" {
		codeAuth.Password = form.Password
	}
	_ , err = db.Handler.XStmt(codeAuth.GetTableName()).Where(dbx.Eq("id",codeAuth.Id)).Cols("user_name","password").Update(&codeAuth)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	data:=map[string]bool{
		"result":true,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}
