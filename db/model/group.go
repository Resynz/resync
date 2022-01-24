/**
 * @Author: Resynz
 * @Date: 2022/1/5 10:35
 */
package model

type Group struct {
	Id   int64  `json:"id"`
	Name string `json:"name" xorm:"varchar(50)"`
}

func (s *Group) GetTableName() string {
	return "group"
}
