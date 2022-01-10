/**
 * @Author: Resynz
 * @Date: 2021/12/31 16:41
 */
package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rosbit/dbx"
	"os"
	"resync/config"
	"resync/db/model"
)

var (
	Handler *dbx.DBI
)

func InitDbHandler() error {
	df:="./data/db"
	f,err := os.Stat(df)
	if err!=nil{
		if os.IsNotExist(err) {
			if err = os.MkdirAll(df, 0755); err != nil {
				return err
			}
		}else {
			return err
		}
	} else {
		if !f.IsDir() {
			return fmt.Errorf("invalid db path:data/db. It should be a dir")
		}
	}
	dbname :="resync.db"
	Handler,err = dbx.CreateDriverDBInstance("sqlite3",fmt.Sprintf("%s/%s",df, dbname),config.Conf.Mode != "release")
	if err!=nil{
		return err
	}
	return nil
}

func InitDBTables() error {
	var err error
	if err = Handler.Sync2(new(model.Admin));err!=nil{
		return err
	}
	if err = Handler.Sync2(new(model.Group));err!=nil{
		return err
	}
	if err = Handler.Sync2(new(model.Task));err!=nil{
		return err
	}
	if err = Handler.Sync2(new(model.TaskDetail));err!=nil{
		return err
	}
	if err = Handler.Sync2(new(model.Action));err!=nil{
		return err
	}
	if err = Handler.Sync2(new(model.CodeAuth));err!=nil{
		return err
	}
	if err = Handler.Sync2(new(model.LoginLog));err!=nil{
		return err
	}
	if err = Handler.Sync2(new(model.TaskLog));err!=nil{
		return err
	}
	// todo init more tables
	return nil
}

func InitDefaultAdmin() error {
	var admin model.Admin
	has,err:=Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("name",config.Conf.DefaultAdminName)).Get(&admin)
	if err!=nil{
		return err
	}
	if !has {
		admin.Name = config.Conf.DefaultAdminName
		admin.Password = config.Conf.DefaultAdminPasswd
		admin.Status = model.AccountStatusEnable
		return Handler.XStmt(admin.GetTableName()).Insert(&admin)
	}
	if admin.Password != config.Conf.DefaultAdminPasswd {
		admin.Password = config.Conf.DefaultAdminPasswd
		_,err = Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id",admin.Id)).Cols("password").Update(&admin)
		return err
	}
	return nil
}
