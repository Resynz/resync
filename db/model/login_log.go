/**
 * @Author: Resynz
 * @Date: 2022/1/5 14:16
 */
package model

type LoginLog struct {
	Id int64 `json:"id"`
	AdminId int64 `json:"admin_id" xorm:"int(11)"`
	Ip string `json:"ip" xorm:"varchar(50)"`
	UserAgent string `json:"user_agent" xorm:"varchar(255)"`
	CreateTime int64 `json:"create_time" xorm:"int(11)"`
}

func (s *LoginLog) GetTableName() string  {
	return "login_log"
}
