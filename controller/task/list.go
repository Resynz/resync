/**
 * @Author: Resynz
 * @Date: 2022/1/5 15:57
 */
package task

import (
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
)

func List(ctx *common.Context) {
	type formValidate struct {
		GroupId int64 `form:"group_id" binding:"" json:"group_id"`
		Page int `form:"page" binding:"" json:"page"`
		Limit int `form:"limit" binding:"" json:"limit"`
	}
	var form formValidate
	if err:=ctx.ShouldBind(&form);err!=nil{
		common.HandleResponse(ctx,code.InvalidParams,nil)
		return
	}
	var task model.Task
	var taskList []*model.Task
	session:=db.Handler.XStmt(task.GetTableName())
	if form.GroupId > 0 {
		session = session.Where(dbx.Eq("group_id",form.GroupId))
	}
	if form.Page > 0 {
		session = session.Limit(form.Limit,(form.Page - 1) * form.Limit)
	}
	err:=session.List(&taskList)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	total,err:=session.Count(&task)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	type listObj struct {
		*model.Task
		CreatorName string `json:"creator_name"`
		ModifierName string `json:"modifier_name"`
		LastRecord struct{
			*model.TaskLog
			CreatorName string `json:"creator_name"`
		} `json:"last_record"`
	}

	list:=make([]*listObj,len(taskList))
	for i,v:=range taskList {
		l:=&listObj{
			Task:         v,
			CreatorName:  "",
			ModifierName: "",
			LastRecord: struct {
				*model.TaskLog
				CreatorName string `json:"creator_name"`
			}{},
		}
		if v.CreatorId > 0 {
			var creator model.Admin
			has,err:=db.Handler.XStmt(creator.GetTableName()).Where(dbx.Eq("id",v.CreatorId)).Get(&creator)
			if err!=nil{
				common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
				return
			}
			if has {
				l.CreatorName = creator.Name
			}
		}
		if v.ModifierId > 0 {
			var creator model.Admin
			has,err:=db.Handler.XStmt(creator.GetTableName()).Where(dbx.Eq("id",v.ModifierId)).Get(&creator)
			if err!=nil{
				common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
				return
			}
			if has {
				l.ModifierName = creator.Name
			}
		}
		var taskLog model.TaskLog
		has,err := db.Handler.XStmt(taskLog.GetTableName()).Where(dbx.Eq("task_id",v.Id)).Desc("id").Get(&taskLog)
		if err!=nil{
			common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
			return
		}

		if has {
			l.LastRecord.TaskLog = &taskLog
			if taskLog.CreatorId > 0 {
				var creator model.Admin
				has,err=db.Handler.XStmt(creator.GetTableName()).Where(dbx.Eq("id",taskLog.CreatorId)).Get(&creator)
				if err!=nil{
					common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
					return
				}
				if has {
					l.LastRecord.CreatorName = creator.Name
				}
			}
		}
		list[i] = l
	}

	data:=map[string]interface{}{
		"list": list,
		"total": total,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}
