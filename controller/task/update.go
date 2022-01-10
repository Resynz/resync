/**
 * @Author: Resynz
 * @Date: 2022/1/5 17:31
 */
package task

import (
	"fmt"
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
)
type UpdateTaskForm struct {
	Note string `form:"note" binding:"" json:"note"`
	SourceCodeType model.SourceCodeType `form:"source_code_type" binding:"" json:"source_code_type"`
	RepositoryUrl string `form:"repository_url" binding:"" json:"repository_url"`
	Branch string `form:"branch" binding:"" json:"branch"`
	CodeAuthId int64 `form:"code_auth_id" binding:"" json:"code_auth_id"`
	ActionList []*model.Action `form:"action_list" binding:"" json:"action_list"`
}
func Update(ctx *common.Context) {

	var form UpdateTaskForm
	_ = ctx.ShouldBind(&form)

	var task model.Task
	has,err:=db.Handler.XStmt(task.GetTableName()).Where(dbx.Eq("id",ctx.Param("id"))).Get(&task)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx,code.InvalidRequest,nil)
		return
	}

	err = db.Handler.Tx(
		dbx.TxStmts(
			updateTaskDetail,
			updateAction,
		),
		dbx.TxArg("task_id",task.Id),
		dbx.TxArg("update_task_form",&form),
	)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	data:=map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}

func updateTaskDetail(stmt *dbx.TxStmt) error {
	form:=stmt.Arg("update_task_form").(*UpdateTaskForm)
	taskId:=stmt.Arg("task_id").(int64)
	var taskDetail model.TaskDetail
	has,err:=stmt.Table(taskDetail.GetTableName()).Where(dbx.Eq("task_id",taskId)).Get(&taskDetail)
	if err!=nil{
		return err
	}
	if !has {
		return fmt.Errorf("任务详情不存在")
	}
	taskDetail.Note = form.Note
	taskDetail.SourceCodeType = form.SourceCodeType
	taskDetail.RepositoryUrl = form.RepositoryUrl
	taskDetail.Branch = form.Branch
	taskDetail.CodeAuthId = form.CodeAuthId
	_ , err = stmt.Table(taskDetail.GetTableName()).Where(dbx.Eq("id",taskDetail.Id)).Cols("note","source_code_type","repository_url","branch","code_auth_id").Update(&taskDetail)
	return err
}

func updateAction(stmt *dbx.TxStmt) error {
	form:=stmt.Arg("update_task_form").(*UpdateTaskForm)
	taskId:=stmt.Arg("task_id").(int64)
	var action model.Action
	if len(form.ActionList) == 0 {
		return stmt.Table(action.GetTableName()).Where(dbx.Eq("task_id",taskId)).Delete(&action)
	}
	ids:=make([]interface{},0)
	for _,v:=range form.ActionList {
		if v.Id > 0 {
			ids = append(ids,v.Id)
		}
	}
	if len(ids) == 0 {
		err:=stmt.Table(action.GetTableName()).Where(dbx.Eq("task_id",taskId)).Delete(&action)
		if err!=nil{
			return err
		}
	}else {
		err:=stmt.Table(action.GetTableName()).Where(dbx.NotIn("id",ids...)).Delete(&action)
		if err!=nil{
			return err
		}
	}
	for _,v:=range form.ActionList {
		if v.Id > 0 {
			has,err:=stmt.Table(action.GetTableName()).Where(dbx.Eq("id",v.Id)).Get(&action)
			if err!=nil{
				return err
			}
			if !has {
				return fmt.Errorf("action not found (id:%d)",v.Id)
			}
			action.Type = v.Type
			action.Content = v.Content
			_,err = stmt.Table(action.GetTableName()).Where(dbx.Eq("id",action.Id)).Cols("type","content").Update(&action)
			if err!=nil{
				return err
			}
			continue
		}
		action = model.Action{
			Id:      0,
			TaskId:  taskId,
			Type:    v.Type,
			Content: v.Content,
		}
		if err := stmt.Table(action.GetTableName()).Insert(&action);err!=nil{
			return err
		}
	}
	return nil
}
