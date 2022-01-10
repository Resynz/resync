/**
 * @Author: Resynz
 * @Date: 2022/1/5 11:41
 */
package model

type Action struct {
	Id int64 `json:"id"`
	TaskId int64 `json:"task_id" xorm:"int(11) not null default 0"`
	Type ActionType `json:"type" xorm:"tinyint(1)"`
	Content string `json:"content" xorm:"text"`
}

func (s *Action) GetTableName() string  {
	return "action"
}
