package config

import (
	"apiTools/utils"
	"errors"
	"gopkg.in/ini.v1"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// 项目配置信息
type appConfig struct {
	serverConfig `conf:"web"`
	redisConfig  `conf:"redis"`
}

// 服务配置信息
type serverConfig struct {
	Port                  string    `conf:"port"`                  // http监听端口
	AppMode               string    `conf:"appMode"`               // app运行模式: production, development
	LogLevel              string    `conf:"logLevel"`              // 日志级别: debug, info, error, warn, panic
	LogSaveDay            uint      `conf:"logSaveDay"`            // 日志文件保留天数
	LogSplitTime          uint      `conf:"logSplitTime"`          // 日志切割时间间隔
	LogOutType            string    `conf:"logOutType"`            // 日志输入类型, json, text
	LogOutPath            string    `conf:"logOutPath"`            // 文件输出位置, file console
	StartTime             time.Time `conf:"startTime"`             // 系统开始运行时间
	EnableIpLimiting      bool      `conf:"enableIpLimiting"`      // 是否开启ip限流
	IpLimitingTimeSeconds uint      `conf:"ipLimitingTimeSeconds"` // IP限流时间段(单位: 秒)
	IpLimitingCount       uint      `conf:"ipLimitingCount"`       // IP限流时间段内请求不能超过的次数
	LiftIpLimiting        uint      `conf:"liftIpLimiting"`        // 解除ip限流的时间(单位: 秒)
}

// redis 配置信息
type redisConfig struct {
	Host     string `conf:"host"`     // redis连接地址
	Port     string `conf:"port"`     // redis连接端口
	Password string `conf:"password"` // redis连接密码
}

var (
	appConf appConfig
)

// 获取指定字段值
// Get("web::port")
// Get("appname")
func Get(ck string) interface{} {
	if ck == "" {
		return nil
	}
	keys := strings.Split(ck, "::")
	if len(keys) == 0 {
		return nil
	}
	appType := reflect.TypeOf(appConf)
	appVal := reflect.ValueOf(appConf)
	for i := 0; i < appType.NumField(); i++ {
		appTypeFiled := appType.Field(i)
		regionTag := appTypeFiled.Tag.Get("conf")
		if regionTag == keys[0] {
			if len(keys) > 1 {
				regionType := appVal.Field(i).Type()
				regionVal := appVal.Field(i)
				for j := 0; j < regionType.NumField(); j++ {
					regionTypeFiled := regionType.Field(j)
					confTag := regionTypeFiled.Tag.Get("conf")
					if confTag == keys[1] {
						val := regionVal.FieldByName(regionTypeFiled.Name).Interface()
						return val
					}
				}
			} else {
				val := appVal.FieldByName(appTypeFiled.Name).Interface()
				return val
			}
		}
	}
	return nil
}

func GetString(ck string) (val string) {
	value := Get(ck)
	if value == nil {
		return
	}
	val, ok := value.(string)
	if !ok {
		return ""
	}
	return
}

func GetInt(ck string) (int) {
	value := Get(ck)
	if value == nil {
		return 0
	}
	if v01, ok01 := value.(uint); ok01 {
		return int(v01)
	}
	if v02, ok02 := value.(int); ok02 {
		return v02
	}
	return 0
}

func GetBool(ck string) (bool) {
	value := Get(ck)
	if value == nil {
		return false
	}
	if v, ok := value.(bool); ok {
		return v
	}

	return false
}

// 初始化配置
func InitConfig() (err error) {
	// 获取配置文件
	configPath := filepath.Join(utils.GetRootPath(), "config", "apitools.ini")
	iniFile, err := ini.Load(configPath)
	if err != nil {
		return
	}
	// 读取web server配置
	err = readServerConfig(iniFile)
	if err != nil {
		return
	}
	// 读取redis配置
	err = readRedisConfig(iniFile)
	if err != nil {
		return
	}
	return
}

// 读取web server配置
func readServerConfig(iniFile *ini.File) (err error) {
	serverConf := iniFile.Section("web")

	httpPort := serverConf.Key("http_port").String()
	if httpPort == "" {
		httpPort = "8091"
	}
	appConf.serverConfig.Port = httpPort

	runMode := serverConf.Key("app_mode").String()
	if runMode == "" {
		runMode = "development"
	}
	appConf.serverConfig.AppMode = runMode

	logLevel := serverConf.Key("logLevel").String()
	if logLevel == "" && runMode == "development" {
		logLevel = "debug"
	} else {
		logLevel = "info"
	}
	appConf.serverConfig.LogLevel = logLevel

	logSaveDay, err := serverConf.Key("logSaveDay").Uint()
	if err != nil {
		logSaveDay = 7
	}
	appConf.serverConfig.LogSaveDay = logSaveDay

	logSplitTime, err := serverConf.Key("logSplitTime").Uint()
	if err != nil {
		logSplitTime = 24
	}
	appConf.serverConfig.LogSplitTime = logSplitTime

	logOutType := serverConf.Key("LogOutType").String()
	if logOutType == "" {
		logOutType = "json"
	}
	appConf.serverConfig.LogOutType = logOutType

	logOutPath := serverConf.Key("logOutPath").String()
	if logOutPath == "" {
		logOutPath = "file"
	}
	appConf.serverConfig.LogOutPath = logOutPath

	startTime := serverConf.Key("startTime").String()
	if startTime == "" {
		appConf.serverConfig.StartTime = time.Now()
	} else {
		runTime, err := time.Parse("2006/01/02", startTime)
		if err != nil {
			appConf.serverConfig.StartTime = time.Now()
		} else {
			appConf.serverConfig.StartTime = runTime
		}
	}

	// ip限流
	enableIpLimiting, err := serverConf.Key("enableIpLimiting").Bool()
	if err != nil {
		enableIpLimiting = false
	}
	appConf.serverConfig.EnableIpLimiting = enableIpLimiting
	if enableIpLimiting {
		ipLimitingTimeSeconds, err := serverConf.Key("ipLimitingTimeSeconds").Uint()
		if err != nil || ipLimitingTimeSeconds == 0 {
			ipLimitingTimeSeconds = 10
		}
		appConf.serverConfig.IpLimitingTimeSeconds = ipLimitingTimeSeconds

		ipLimitingCount, err := serverConf.Key("ipLimitingCount").Uint()
		if err != nil || ipLimitingCount == 0 {
			ipLimitingCount = 8
		}
		appConf.serverConfig.IpLimitingCount = ipLimitingCount

		liftIpLimiting, err := serverConf.Key("liftIpLimiting").Uint()
		if err != nil || liftIpLimiting == 0 {
			liftIpLimiting = 5
		}
		appConf.serverConfig.LiftIpLimiting = liftIpLimiting
	}

	return
}

// 读取redis配置
func readRedisConfig(iniFile *ini.File) (err error) {
	redisConf := iniFile.Section("redis")

	host := redisConf.Key("host").String()
	if host == "" {
		return errors.New("config file redis host can not be empty")
	}
	appConf.redisConfig.Host = host

	port := redisConf.Key("port").String()
	if port == "" {
		port = "6379"
	}
	appConf.redisConfig.Port = port

	password := redisConf.Key("password").String()
	appConf.redisConfig.Password = password

	return
}
