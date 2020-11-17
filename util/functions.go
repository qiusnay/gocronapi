package util

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

var MyAppConfig *MyConfig

type MyConfig struct {
	IgnoreURLs []interface{}
	JWTTimeout int
	LogLevel   string
	Secret     string
}

const (
	// key定义
	CODE           string = "code"
	MSG            string = "msg"
	DATA           string = "data"
	LoginSuccess   string = "恭喜, 登录成功"
	LoginFailur    string = "登录失败"
	UsernameFailur string = "用户名错误"
	PasswordFailur string = "密码错误"
)

var BashPath string

func SetBASEPath() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	BashPath = strings.Replace(dir, "\\", "/", -1)
}

//载入配置
func LoadConf() iris.Configuration {
	myConf := iris.YAML("./conf/app.yml")
	conf := MyConfig{}
	conf.IgnoreURLs = myConf.Other["IgnoreURLs"].([]interface{})
	conf.JWTTimeout = myConf.Other["JWTTimeout"].(int)
	conf.LogLevel = myConf.Other["LogLevel"].(string)
	conf.Secret = myConf.Other["Secret"].(string)
	MyAppConfig = &conf
	// fmt.Println(fmt.Sprintf("init app %+v", myConf))
	return myConf
}

type App struct {
	DisablePathCorrection             bool      `yaml:"DisablePathCorrection"`
	EnablePathEscape                  bool      `yaml:"EnablePathEscape"`
	FireMethodNotAllowed              bool      `yaml:"FireMethodNotAllowed"`
	DisableBodyConsumptionOnUnmarshal bool      `yaml:"DisableBodyConsumptionOnUnmarshal"`
	TimeFormat                        time.Time `yaml:"TimeFormat"`
	Charset                           string    `yaml:"Charset"`
}

//两个结构体反射覆盖
// func GetMyConf() iris.Configuration {
// 	Log = NewLogger()
// 	myConf := iris.YAML("./conf/app.yml")
// 	// Log.Info(fmt.Sprintf("%+v", myConf.App))
// 	irisDefaultConfig := iris.DefaultConfiguration()
// 	irisConfVal := reflect.ValueOf(&irisDefaultConfig)
// 	irisConfKey := reflect.TypeOf(irisDefaultConfig)
// 	for irisk := 0; irisk < irisConfKey.NumField(); irisk++ {
// 		irisKey := irisConfKey.Field(irisk).Name // a reflect.StructField
// 		v := reflect.ValueOf(myConf.App)
// 		t := reflect.TypeOf(myConf.App)
// 		for k := 0; k < t.NumField(); k++ {
// 			myKey := t.Field(k).Name
// 			if myKey == irisKey {
// 				myval := v.Field(k).Interface()
// 				switch value := myval.(type) {
// 				case string:
// 					irisConfVal.Elem().FieldByName(irisKey).SetString(value)
// 				case bool:
// 					irisConfVal.Elem().FieldByName(irisKey).SetBool(value)
// 				case int64:
// 					irisConfVal.Elem().FieldByName(irisKey).SetInt(value)
// 				}
// 			}
// 		}
// 	}
// 	// Log.Info(fmt.Sprintf("%+v", irisDefaultConfig))
// 	return irisDefaultConfig
// }

func JsonSuccess(ctx iris.Context, msg string, data interface{}) {
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		CODE: iris.StatusOK,
		MSG:  msg,
		DATA: data,
	})
}

// 401 error define
func Unauthorized(ctx iris.Context, msg string, data interface{}) {
	unauthorized := iris.StatusUnauthorized

	ctx.StatusCode(unauthorized)
	ctx.JSON(iris.Map{
		CODE: unauthorized,
		MSG:  msg,
		DATA: data,
	})
}

// common error define
func JsonFail(ctx iris.Context, status int, msg string, data interface{}) {
	ctx.StatusCode(status)
	ctx.JSON(iris.Map{
		CODE: status,
		MSG:  msg,
		DATA: data,
	})
}

func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}
	pathMap := lfshook.PathMap{
		logrus.InfoLevel: "./log/cronlog." + time.Now().Format("2006-01-02") + ".log",
	}
	Log = logrus.New()
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	return Log
}

func InArray(need interface{}, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}
