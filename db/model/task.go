/**
 * @Author: Resynz
 * @Date: 2022/1/5 10:39
 */
package model

type Task struct {
	Id int64 `json:"id"`
	Name string `json:"name" xorm:"varchar(50)"`
	GroupId int64 `json:"group_id" xorm:"int(11) not null default 0"`
	CreatorId int64 `json:"creator_id" xorm:"int(11) not null default 0"`
	CreateTime int64 `json:"create_time" xorm:"int(11) not null default 0"`
	ModifierId int64 `json:"modifier_id" xorm:"int(11) not null default 0"`
	ModifyTime int64 `json:"modify_time" xorm:"int(11) not null default 0"`
}

func (s *Task) GetTableName() string  {
	return "task"
}
