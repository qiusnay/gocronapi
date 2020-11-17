package api

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/qiusnay/gocronapi/middleware/jwt"
	"github.com/qiusnay/gocronapi/models"
	"github.com/qiusnay/gocronapi/models/db"
	"github.com/qiusnay/gocronapi/util"
)

type User struct{}

//获取用户信息
func (u *User) GetUserinfo(ctx iris.Context) {
	token, _ := jwt.Myjwt.FromAuthHeader(ctx)
	Userinfo, err := jwt.Myjwt.DecodeToken(ctx, token)
	if err != nil {
		util.JsonFail(ctx, iris.StatusInternalServerError, "token不正确", nil)
	}
	// fmt.Println(fmt.Sprintf("token data : %+v", token))
	// fmt.Println(fmt.Sprintf("user data : %+v", Userinfo))
	mUser := models.GetUserByUserBUserId(Userinfo.Id)
	userMap := make(map[string]interface{})
	userMap["userid"] = mUser[0].User.Id
	userMap["username"] = mUser[0].User.Username
	userMap["lastlogin"] = mUser[0].User.LastLogin
	userMap["email"] = mUser[0].User.Email
	userMap["user_code"] = mUser[0].User.UserCode
	userMap["is_active"] = mUser[0].User.IsActive

	userMap["roles"] = strings.Split(mUser[0].UserGroup.Roles, ",")

	// fmt.Println(fmt.Sprintf("db data : %+v", mUser))
	// if dberr != nil {
	// 	ctx.Application().Logger().Errorf("用户[%s]获取信息失败。%s", mUser.Username, dberr.Error())
	// 	util.JsonFail(ctx, iris.StatusInternalServerError, util.LoginFailur, nil)
	// }
	//密码不能传到前端去
	// mUser.Password = ""
	util.JsonSuccess(ctx, "success", userMap)
}

//用户登录
func (u *User) PostLogin(ctx iris.Context) {
	user := new(db.User)
	if err := ctx.ReadJSON(&user); err != nil {
		ctx.Application().Logger().Errorf("用户[%s]登录失败。%s", "", err.Error())
		util.JsonFail(ctx, iris.StatusBadRequest, util.LoginFailur, nil)
		return
	}

	mUser := new(db.User)
	mUser.Username = user.Username
	has, err := models.GetUserByUserByUserName(mUser)
	if err != nil {
		ctx.Application().Logger().Errorf("用户[%s]登录失败。%s", user.Username, err.Error())
		util.JsonFail(ctx, iris.StatusInternalServerError, util.LoginFailur, nil)
		return
	}
	if !has { // 用户名不正确
		util.Unauthorized(ctx, util.UsernameFailur, nil)
		return
	}
	ckPassword := util.Verify(user.Password, mUser.Password)
	if !ckPassword {
		util.Unauthorized(ctx, util.PasswordFailur, nil)
		return
	}
	token, err := jwt.GenerateToken(mUser)
	// golog.Infof("用户[%s], 登录生成token [%s]", mUser.Username, token)
	resultData := map[string]string{"token": token}
	util.JsonSuccess(ctx, "success", resultData)
}

func (u *User) PostLogout(ctx iris.Context) {

	// util.JsonSuccess(ctx, "success", resultData)
}
