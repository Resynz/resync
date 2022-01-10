/**
 * @Author: Resynz
 * @Date: 2022/1/5 15:37
 */
package group

import (
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
)

func List(ctx *common.Context) {
	var group model.Group
	var list []*model.Group
	if err:=db.Handler.XStmt(group.GetTableName()).List(&list);err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	data:=map[string]interface{}{
		"list": list,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}

func Add(ctx *common.Context) {
	type formValidate struct {
		Name string `form:"name" binding:"required" json:"name"`
	}
	var form formValidate
	if err:=ctx.ShouldBind(&form);err!=nil{
		common.HandleResponse(ctx,code.InvalidParams,nil)
		return
	}
	group:=model.Group{
		Id:   0,
		Name: form.Name,
	}
	if err:=db.Handler.XStmt(group.GetTableName()).Insert(&group);err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	data:=map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}

func Update(ctx *common.Context) {
	var group model.Group
	has,err:=db.Handler.XStmt(group.GetTableName()).Where(dbx.Eq("id",ctx.Param("id"))).Get(&group)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx,code.InvalidRequest,nil)
		return
	}
	type formValidate struct {
		Name string `form:"name" binding:"required" json:"name"`
	}
	var form formValidate
	if err=ctx.ShouldBind(&form);err!=nil{
		common.HandleResponse(ctx,code.InvalidParams,nil)
		return
	}
	group.Name = form.Name
	if _,err = db.Handler.XStmt(group.GetTableName()).Where(dbx.Eq("id",group.Id)).Cols("name").Update(&group);err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	data:=map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}

func Delete(ctx *common.Context) {
	var group model.Group
	has,err:=db.Handler.XStmt(group.GetTableName()).Where(dbx.Eq("id",ctx.Param("id"))).Get(&group)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx,code.InvalidRequest,nil)
		return
	}
	if err = db.Handler.XStmt(group.GetTableName()).Where(dbx.Eq("id",group.Id)).Delete(&group);err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	data:=map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}
