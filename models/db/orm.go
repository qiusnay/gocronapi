package db

import "time"

// mysql user table
type User struct {
	Id         int64     `xorm:"pk autoincr INT(10) notnull" json:"id" form:"id" `
	Username   string    `xorm:"VARCHAR(50) notnull comment('用户名') index(IX_username)"  json:"username" form:"username"`
	Password   string    `xorm:"VARCHAR(100) notnull comment('密码')" json:"password" form:"password"`
	LastLogin  string    `xorm:"VARCHAR(50) notnull comment('上次登录时间')" json:"name" form:"name"`
	UserCode   string    `xorm:"int notnull comment('用户角色')" json:"phone" form:"phone"`
	Email      string    `xorm:"VARCHAR(50) notnull comment('邮箱')" json:"email" form:"email"`
	IsActive   string    `xorm:"int notnull comment('是否激活')" json:"userface" form:"userface"`
	CreateTime time.Time `xorm:"timestamp notnull default CURRENT_TIMESTAMP comment('创建时间')" json:"VARCHAR(50) createTime" form:"createTime"`
	UpdateTime time.Time `xorm:"timestamp notnull default CURRENT_TIMESTAMP comment('更新时间')" json:"VARCHAR(50) updateTime" form:"updateTime" `
}

func (t *User) TableName() string {
	return "tb_user"
}

type UserGroup struct {
	Id         int64     `xorm:"pk autoincr INT(10) notnull" json:"id" form:"id" `
	Userid     string    `xorm:"VARCHAR(50) notnull comment('用户ID') index(IX_userid)"  json:"userid" form:"userid"`
	Roles      string    `xorm:"VARCHAR(100) notnull comment('权限列表')" json:"roles" form:"roles"`
	CreateTime time.Time `xorm:"timestamp notnull default CURRENT_TIMESTAMP comment('创建时间')" json:"VARCHAR(50) createTime" form:"createTime"`
	UpdateTime time.Time `xorm:"timestamp notnull default CURRENT_TIMESTAMP comment('更新时间')" json:"VARCHAR(50) updateTime" form:"updateTime" `
}

func (t *UserGroup) TableName() string {
	return "tb_user_roles"
}

type UserInfoAll struct {
	User      `xorm:"extends"`
	UserGroup `xorm:"extends"`
}
