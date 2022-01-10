/**
 * @Author: Resynz
 * @Date: 2022/1/5 10:52
 */
package model

type TaskDetail struct {
	Id int64 `json:"id"`
	TaskId int64 `json:"task_id" xorm:"int(11) not null default 0"`
	Note string `json:"note" xorm:"text"`
	SourceCodeType SourceCodeType `json:"source_code_type" xorm:"tinyint(1) not null default 0"`
	RepositoryUrl string `json:"repository_url" xorm:"varchar(255)"`
	Branch string `json:"branch" xorm:"varchar(255)"`
	CodeAuthId int64 `json:"code_auth_id" xorm:"int(11) not null default 0"`
}

func (s *TaskDetail) GetTableName() string {
	return "task_detail"
}
