/**
 * @Author: Resynz
 * @Date: 2022/1/7 11:05
 */
package code_auth

import (
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
)

func Delete(ctx *common.Context) {
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
	if err = db.Handler.XStmt(codeAuth.GetTableName()).Where(dbx.Eq("id",codeAuth.Id)).Delete(&codeAuth);err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	data:=map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}
