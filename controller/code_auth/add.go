/**
 * @Author: Resynz
 * @Date: 2022/1/7 10:58
 */
package code_auth

import (
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
	"time"
)

func Add(ctx *common.Context) {
	type formValidate struct {
		AuthType model.AuthType `form:"auth_type" binding:"required" json:"auth_type"`
		UserName string `form:"user_name" binding:"required" json:"user_name"`
		Password string `form:"password" binding:"required" json:"password"`
	}
	var form formValidate
	if err:=ctx.ShouldBind(&form);err!=nil{
		common.HandleResponse(ctx,code.InvalidParams,nil)
		return
	}
	codeAuth:=model.CodeAuth{
		Id:         0,
		AuthType:   form.AuthType,
		UserName:   form.UserName,
		Password:   form.Password,
		CreatorId:  ctx.Auth.Id,
		CreateTime: time.Now().Unix(),
		ModifierId: ctx.Auth.Id,
		ModifyTime: time.Now().Unix(),
	}
	if err:=db.Handler.XStmt(codeAuth.GetTableName()).Insert(&codeAuth);err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	data:=map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}
