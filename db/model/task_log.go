/**
 * @Author: Resynz
 * @Date: 2022/1/5 16:03
 */
package model

type TaskLog struct {
	Id int64 `json:"id"`
	TaskId int64 `json:"task_id" xorm:"int(11)"`
	Status TaskStatus `json:"status" xorm:"tinyint(1)"`
	StartTime int64 `json:"start_time" xorm:"int(11)"`
	EndTime int64 `json:"end_time" xorm:"int(11)"`
	CreatorId int64 `json:"creator_id" xorm:"int(11)"`
}

func (s *TaskLog) GetTableName() string  {
	return "task_log"
}
