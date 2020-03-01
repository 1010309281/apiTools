package modles

import (
	"apiTools/libs/config"
	"apiTools/utils"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"path/filepath"
	"time"
)

var (
	RedisPool *redis.Pool
	JsonData  interface{}
	SqlConn   *gorm.DB
)

func InitRedis() (err error) {
	RedisPool = &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%s",
				config.GetString("redis::host"), config.GetString("redis::port")))
			if err != nil {
				return nil, err
			}
			if config.GetString("redis::password") != "" {
				if _, err := conn.Do("AUTH", config.GetString("redis::password")); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return
		},
		MaxIdle:     16,              // 最大空闲连接数
		MaxActive:   32,              // 最大活跃连接数
		IdleTimeout: time.Second * 3, // 最大空闲超时时间
	}
	pool := RedisPool.Get()
	defer pool.Close()
	if err := pool.Err(); err != nil {
		return fmt.Errorf("init redis failed, err: %v", err)
	}
	return
}

func InitMysql() (err error) {
	// CREATE DATABASE `apitools` CHARACTER SET utf8 COLLATE utf8_general_ci;
	// ALTER user 'apitools'@'localhost' IDENTIFIED BY 'apitools#.*'


	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.GetString("mysql::user"),
			config.GetString("mysql::password"),
			config.GetString("mysql::host"),
			config.GetString("mysql::port"),
			config.GetString("mysql::db"),
		))
	if err != nil {
		return
	}
	SqlConn = db

	if config.GetBool("mysql::enableDebug") {
		db.LogMode(true)
	}
	// 单数表名
	SqlConn.SingularTable(true)

	// 自动映表
	SqlConn.AutoMigrate(&ProxyPool{})

	return
}

// 关闭IO流
func CloseIO() {
	defer func() {
		recover()
	}()

	// 关闭redis io
	RedisPool.Close()

	// 关闭mysql io
	SqlConn.Close()
}

// 初始化api docs json数据
func InitApiDocsJsonData() (err error) {
	jsonFile, err := ioutil.ReadFile(filepath.Join(utils.GetRootPath(), "data/api", "apidocs.json"))
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonFile, &JsonData)
	if err != nil {
		return
	}
	return
}
