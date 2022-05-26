/**
 * @Author: Resynz
 * @Date: 2022/1/5 14:19
 */
package admin

import (
	"fmt"
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
)

// Info 获取管理员信息
func Info(ctx *common.Context) {
	data := map[string]string{
		"name": ctx.Auth.Name,
	}
	common.HandleResponse(ctx, code.SuccessCode, data)
}

// List 获取管理员列表
func List(ctx *common.Context) {
	type formValidate struct {
		Name   string              `form:"name" binding:"" json:"name"`
		Page   int                 `form:"page" binding:"" json:"page"`
		Limit  int                 `form:"limit" binding:"" json:"limit"`
		Status model.AccountStatus `form:"status" binding:"" json:"status"`
	}
	var form formValidate
	_ = ctx.ShouldBind(&form)
	var admin model.Admin
	var adminList []*model.Admin
	session := db.Handler.XStmt(admin.GetTableName())
	if form.Name != "" {
		session = session.Where(dbx.Op("name", "like", fmt.Sprintf("%%%s%%", form.Name)))
	}
	if form.Status != model.AccountStatusUnknown {
		session = session.Where(dbx.Eq("status", form.Status))
	}
	if form.Page > 0 {
		session = session.Limit(form.Limit, (form.Page-1)*form.Limit)
	}
	err := session.List(&adminList)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	total, err := session.Count(&admin)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	type listObj struct {
		Id     int64               `json:"id"`
		Name   string              `json:"name"`
		Status model.AccountStatus `json:"status"`
	}
	list := make([]*listObj, len(adminList))
	for i, v := range adminList {
		l := &listObj{
			Id:     v.Id,
			Name:   v.Name,
			Status: v.Status,
		}
		list[i] = l
	}
	data := map[string]interface{}{
		"total": total,
		"list":  list,
	}
	common.HandleResponse(ctx, code.SuccessCode, data)
}

// Disable 禁用管理员
func Disable(ctx *common.Context) {
	type formValidate struct {
		Id int64 `form:"id" binding:"required" json:"id"`
	}
	var form formValidate
	if err := ctx.ShouldBind(&form); err != nil {
		common.HandleResponse(ctx, code.InvalidParams, nil)
		return
	}
	var admin model.Admin
	has, err := db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id", form.Id)).Get(&admin)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}
	admin.Status = model.AccountStatusDisable
	if _, err = db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id", admin.Id)).Cols("status").Update(&admin); err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	data := map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx, code.SuccessCode, data)
}

// Enable 启用管理员
func Enable(ctx *common.Context) {
	type formValidate struct {
		Id int64 `form:"id" binding:"required" json:"id"`
	}
	var form formValidate
	if err := ctx.ShouldBind(&form); err != nil {
		common.HandleResponse(ctx, code.InvalidParams, nil)
		return
	}
	var admin model.Admin
	has, err := db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id", form.Id)).Get(&admin)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}
	admin.Status = model.AccountStatusEnable
	if _, err = db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id", admin.Id)).Cols("status").Update(&admin); err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	data := map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx, code.SuccessCode, data)
}

// Add 新增管理员
func Add(ctx *common.Context) {
	type formValidate struct {
		Name     string `form:"name" binding:"required" json:"name"`
		Password string `form:"password" binding:"required" json:"password"`
	}
	var form formValidate
	if err := ctx.ShouldBind(&form); err != nil {
		common.HandleResponse(ctx, code.InvalidParams, nil)
		return
	}
	var admin model.Admin
	has, err := db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("name", form.Name)).Get(&admin)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	if has {
		common.HandleResponse(ctx, code.InvalidRequest, nil, "该账号已存在")
		return
	}

	admin = model.Admin{
		Id:       0,
		Name:     form.Name,
		Password: form.Password,
		Status:   model.AccountStatusEnable,
	}

	if err = db.Handler.XStmt(admin.GetTableName()).Insert(&admin); err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	data := map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx, code.SuccessCode, data)
}

// Update 修改管理员
func Update(ctx *common.Context) {
	var admin model.Admin
	has, err := db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id", ctx.Param("id"))).Get(&admin)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}
	type formValidate struct {
		Name     string `form:"name" binding:"required" json:"name"`
		Password string `form:"password" binding:"" json:"password"`
	}
	var form formValidate
	if err = ctx.ShouldBind(&form); err != nil {
		common.HandleResponse(ctx, code.InvalidParams, nil)
		return
	}
	var check model.Admin
	has, err = db.Handler.XStmt(check.GetTableName()).Where(dbx.Eq("name", form.Name), dbx.Op("id", "!=", admin.Id)).Get(&check)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	if has {
		common.HandleResponse(ctx, code.InvalidRequest, nil, "该账号已存在")
		return
	}
	admin.Name = form.Name
	if form.Password != "" {
		admin.Password = form.Password
	}
	if _, err = db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id", admin.Id)).Cols("name", "password").Update(&admin); err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	data := map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx, code.SuccessCode, data)
}

// Delete 删除管理员
func Delete(ctx *common.Context) {
	var admin model.Admin
	has, err := db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id", ctx.Param("id"))).Get(&admin)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}
	if err = db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id", admin.Id)).Delete(&admin); err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	data := map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx, code.SuccessCode, data)
}

// ModifyPasswd 修改密码
func ModifyPasswd(ctx *common.Context) {
	var form struct {
		OldPass string `form:"old_pass" binding:"required" json:"old_pass"`
		NewPass string `form:"new_pass" binding:"required" json:"new_pass"`
	}
	if err := ctx.ShouldBind(&form); err != nil {
		common.HandleResponse(ctx, code.InvalidParams, nil)
		return
	}

	var admin model.Admin
	_, err := db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id", ctx.Auth.Id)).Get(&admin)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	if admin.Password != form.OldPass {
		common.HandleResponse(ctx, code.InvalidRequest, nil, "原密码有误")
		return
	}
	admin.Password = form.NewPass
	_, err = db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id", admin.Id)).Cols("password").Update(&admin)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	data := map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx, code.SuccessCode, data)
}
