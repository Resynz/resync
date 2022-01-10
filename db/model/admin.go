/**
 * @Author: Resynz
 * @Date: 2021/12/31 17:07
 */
package model


type Admin struct {
	Id int64 `json:"id"`
	Name string `json:"name" xorm:"varchar(50)"`
	Password string `json:"password" xorm:"varchar(255)"`
	Status AccountStatus `json:"status" xorm:"tinyint(1)"`
}

func (s *Admin) GetTableName() string  {
	return "admin"
}