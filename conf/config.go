package conf

import (
	"github.com/go-ini/ini"
	"log"
	"time"
	_ "time"
)

type App struct {
	Name      string
	Host      string
	Url       string
	JwtSecret string
	PwdSalt   string
}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	DbType string
	DbUser string
	DbPass string
	DbHost string
	DbPort int
	DbName string
}

var AppSetting = &App{}
var ServerSetting = &Server{}
var DatabaseSetting = &Database{}
var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app_debug.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app_debug.ini': %v", err)
	}
	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}
