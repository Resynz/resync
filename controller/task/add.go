/**
 * @Author: Resynz
 * @Date: 2022/1/5 16:25
 */
package task

import (
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
	"time"
)

func Add(ctx *common.Context) {
	type formValidate struct {
		GroupId int64 `form:"group_id" binding:"" json:"group_id"`
		Name string `form:"name" binding:"required" json:"name"`
	}
	var form formValidate
	if err:=ctx.ShouldBind(&form);err!=nil{
		common.HandleResponse(ctx,code.InvalidParams,nil,err.Error())
		return
	}
	task:=&model.Task{
		Name:       form.Name,
		GroupId:    form.GroupId,
		CreatorId:  ctx.Auth.Id,
		CreateTime: time.Now().Unix(),
		ModifierId: ctx.Auth.Id,
		ModifyTime: time.Now().Unix(),
	}
	if err := db.Handler.Tx(
		dbx.TxStmts(
			createTask,
		),
		dbx.TxArg("task",task),
	);err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}

	data:=map[string]int64 {
		"task_id": task.Id,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}

func createTask(stmt *dbx.TxStmt) error {
	task:=stmt.Arg("task").(*model.Task)
	if err:=stmt.Table(task.GetTableName()).Insert(task);err!=nil{
		return err
	}
	taskDetail:=model.TaskDetail{
		Id:             0,
		TaskId:         task.Id,
		Note:           "",
		SourceCodeType: model.SourceCodeTypeNone,
		RepositoryUrl:  "",
		Branch:         "",
		CodeAuthId:     0,
	}
	return stmt.Table(taskDetail.GetTableName()).Insert(&taskDetail)
}
