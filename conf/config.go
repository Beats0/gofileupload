package conf

import (
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"time"
	_ "time"
)

type App struct {
	Name           string
	Url            string
	PwdSalt        string
	JwtSecret      string
	CheckSecretKey string
}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxDiskSize  int
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
	var iniConfig string
	var mode = flag.String("mode", "", "")
	flag.Parse()
	if *mode == "release" {
		iniConfig = "conf/app_release.ini"
	} else {
		iniConfig = "conf/app_debug.ini"
	}
	fmt.Println("Load config:", iniConfig)
	cfg, err = ini.Load(iniConfig)
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
