package util

import (
	"os"
	"reflect"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Log *logrus.Logger

type Config struct {
	App App `yaml:"App"`
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
func GetMyConf() iris.Configuration {
	Log = NewLogger()
	myConf, err := ReadYamlConfig("./conf/app.yml")
	if err != nil {
		Log.Error("Error. %s", err)
	}
	// Log.Info(fmt.Sprintf("%+v", myConf.App))
	irisDefaultConfig := iris.DefaultConfiguration()
	irisConfVal := reflect.ValueOf(&irisDefaultConfig)
	irisConfKey := reflect.TypeOf(irisDefaultConfig)
	for irisk := 0; irisk < irisConfKey.NumField(); irisk++ {
		irisKey := irisConfKey.Field(irisk).Name // a reflect.StructField
		v := reflect.ValueOf(myConf.App)
		t := reflect.TypeOf(myConf.App)
		for k := 0; k < t.NumField(); k++ {
			myKey := t.Field(k).Name
			if myKey == irisKey {
				myval := v.Field(k).Interface()
				switch value := myval.(type) {
				case string:
					irisConfVal.Elem().FieldByName(irisKey).SetString(value)
				case bool:
					irisConfVal.Elem().FieldByName(irisKey).SetBool(value)
				case int64:
					irisConfVal.Elem().FieldByName(irisKey).SetInt(value)
				}
			}
		}
	}
	// Log.Info(fmt.Sprintf("%+v", irisDefaultConfig))
	return irisDefaultConfig
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

//read yaml config
//注：path为yaml或yml文件的路径
func ReadYamlConfig(path string) (*Config, error) {
	conf := &Config{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		yaml.NewDecoder(f).Decode(conf)
	}
	return conf, nil
}

func InArray(need interface{}, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}
