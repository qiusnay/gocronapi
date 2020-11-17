package models

import (
	"github.com/qiusnay/gocronapi/models/db"
)

func GetUserByUserByUserName(user *db.User) (bool, error) {
	e := db.MasterEngine()
	return e.Get(user)
}

func GetUserByUserBUserId(userid int64) []db.UserInfoAll {
	e := db.MasterEngine()
	users := make([]db.UserInfoAll, 0)
	e.Table("tb_user").Join("INNER", "tb_user_roles", "tb_user_roles.userid = tb_user.id").Where("tb_user.id = ?", userid).Find(&users)
	return users
	// e := db.MasterEngine()
	// users := make([]db.UserInfoAll, 0)
	// return e.Alias("u").Join("INNER", "tb_user_roles s", "s.userid=u.id").Get(user)
}
