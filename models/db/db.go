package db

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
	"xorm.io/core"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/kataras/golog"
	"github.com/qiusnay/gocronapi/util"
)

var (
	masterEngine *xorm.Engine
	slaveEngine  *xorm.Engine
	lock         sync.Mutex
)

type DB struct {
	Master DBConfigInfo
	Slave  DBConfigInfo
}

type DBConfigInfo struct {
	Dialect      string `yaml:"dialect"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Database     string `yaml:"database"`
	Charset      string `yaml:"charset"`
	ShowSql      bool   `yaml:"showSql"`
	LogLevel     string `yaml:"logLevel"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
}

func LoadDbConf(dbtype string) DBConfigInfo {
	conf := new(DB)
	yamlFile, _ := ioutil.ReadFile(util.BashPath + "/conf/db.yml")
	yaml.Unmarshal(yamlFile, conf)
	if dbtype == "master" {
		return conf.Master
	} else {
		return conf.Slave
	}
}

// 主库，单例
func MasterEngine() *xorm.Engine {
	if masterEngine != nil {
		return masterEngine
	}

	lock.Lock()
	defer lock.Unlock()

	if masterEngine != nil {
		return masterEngine
	}

	master := LoadDbConf("master")
	engine, err := xorm.NewEngine(master.Dialect, GetConnURL(&master))
	if err != nil {
		fmt.Println(fmt.Sprintf("db connect error : %+v", err))
		return nil
	}
	migrateDb(engine)
	engine.SetMapper(core.GonicMapper{})

	engine.ShowSQL(master.ShowSql)
	SysTimeLocation, _ := time.LoadLocation("Asia/Chongqing") // 中国时区
	engine.SetTZLocation(SysTimeLocation)
	if master.MaxIdleConns > 0 {
		engine.SetMaxIdleConns(master.MaxIdleConns)
	}
	if master.MaxOpenConns > 0 {
		engine.SetMaxOpenConns(master.MaxOpenConns)
	}

	masterEngine = engine
	return masterEngine
}

// 从库，单例
func SlaveEngine() *xorm.Engine {
	if slaveEngine != nil {
		return slaveEngine
	}

	lock.Lock()
	defer lock.Unlock()

	if slaveEngine != nil {
		return slaveEngine
	}

	slave := LoadDbConf("slave")
	engine, err := xorm.NewEngine(slave.Dialect, GetConnURL(&slave))
	if err != nil {
		golog.Fatalf("@@@ Instance Slave DB error!! %s", err)
		return nil
	}

	engine.ShowSQL(slave.ShowSql)
	SysTimeLocation, _ := time.LoadLocation("Asia/Chongqing") // 中国时区
	engine.SetTZLocation(SysTimeLocation)
	if slave.MaxIdleConns > 0 {
		engine.SetMaxIdleConns(slave.MaxIdleConns)
	}
	if slave.MaxOpenConns > 0 {
		engine.SetMaxOpenConns(slave.MaxOpenConns)
	}

	slaveEngine = engine
	return engine
}

//创建表结构
func migrateDb(engine *xorm.Engine) {
	exist, _ := engine.IsTableExist(&User{})
	if !exist {
		engine.CreateTables(&User{})
		engine.CreateIndexes(&UserGroup{})

		engine.CreateTables(&UserGroup{})
		engine.CreateIndexes(&UserGroup{})
	}
}

// 获取数据库连接的url
// true：master主库
func GetConnURL(info *DBConfigInfo) (url string) {
	//db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		info.User,
		info.Password,
		info.Host,
		info.Port,
		info.Database,
		info.Charset)
	//golog.Infof("@@@ DB conn==>> %s", url)
	return
}
