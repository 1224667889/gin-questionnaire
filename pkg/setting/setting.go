package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	User        string
	Port		int
	Password    string
	Host        string
	Name        string
}

var DatabaseSetting = &Database{}

type Mongo struct {
	Host 	string
	Port 	int
	Name 	string
}

var MongoSetting = &Mongo{}

var cfg *ini.File

// Setup APP基础设置
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("无法解析配置文件'conf/app.ini': %v", err)
	}

	mapSetting("app", AppSetting)
	mapSetting("server", ServerSetting)
	mapSetting("database", DatabaseSetting)
	mapSetting("mongo", MongoSetting)

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
}

// mapSetting 载入配置
func mapSetting(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("配置[%s]载入错误: %v", section, err)
	}
}
