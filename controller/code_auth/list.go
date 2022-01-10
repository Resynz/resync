/**
 * @Author: Resynz
 * @Date: 2022/1/7 10:53
 */
package code_auth

import (
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
)

func List(ctx *common.Context) {
	type formValidate struct {
		AuthType model.AuthType `form:"auth_type" binding:"required" json:"auth_type"`
	}
	var form formValidate
	if err:=ctx.ShouldBind(&form);err!=nil{
		common.HandleResponse(ctx,code.InvalidParams,nil)
		return
	}
	var codeAuth model.CodeAuth
	var codeAuthList []*model.CodeAuth
	err:=db.Handler.XStmt(codeAuth.GetTableName()).Where(dbx.Eq("auth_type",form.AuthType)).List(&codeAuthList)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	type listObj struct {
		Id int64 `json:"id"`
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	list:=make([]*listObj,len(codeAuthList))
	for i,v:=range codeAuthList {
		l:=&listObj{
			Id:       v.Id,
			UserName: v.UserName,
			Password: "***",
		}
		list[i] = l
	}
	data:=map[string]interface{}{
		"list": list,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}
