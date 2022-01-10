/**
 * @Author: Resynz
 * @Date: 2022/1/5 11:13
 */
package model

type CodeAuth struct {
	Id int64 `json:"id"`
	AuthType AuthType `json:"auth_type" xorm:"tinyint(1)"`
	UserName string `json:"user_name" xorm:"varchar(100)"`
	Password string `json:"password" xorm:"varchar(255)"`
	CreatorId int64 `json:"creator_id" xorm:"int(11) not null default 0"`
	CreateTime int64 `json:"create_time" xorm:"int(11) not null default 0"`
	ModifierId int64 `json:"modifier_id" xorm:"int(11) not null default 0"`
	ModifyTime int64 `json:"modify_time" xorm:"int(11) not null default 0"`
}

func (s *CodeAuth) GetTableName() string  {
	return "code_auth"
}
