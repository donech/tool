package xdb

import (
	"context"
	"testing"

	"github.com/donech/tool/xtrace"

	"github.com/donech/tool/xlog"
)

type BaseApps struct {
	AppId        string `json:"app_id"`
	AppName      string `json:"app_name"`
	DebugMode    int    `json:"debug_mode"`
	AppConfig    string `json:"app_config"`
	Status       string `json:"status"`
	Webpath      string `json:"webpath"`
	Description  string `json:"description"`
	LocalVer     string `json:"local_ver"`
	RemoteVer    string `json:"remote_ver"`
	AuthorName   string `json:"author_name"`
	AuthorUrl    string `json:"author_url"`
	AuthorEmail  string `json:"author_email"`
	Dbver        string `json:"dbver"`
	RemoteConfig string `json:"remote_config"`
}

func TestOpen(t *testing.T) {
	logConfig := xlog.Config{
		ServiceName: "db_test",
		Level:       "debug",
		LevelColor:  false,
		Format:      "json",
		Stdout:      false,
		File: xlog.FileLogConfig{
			Filename:  "test.log",
			LogRotate: false,
		},
	}
	_, err := xlog.New(logConfig)
	if err != nil {
		t.Error("init logger failed")
	}

	dbConfig := Config{
		Dsn:         "root:example@tcp(localhost:3307)/b2b2c_dev?charset=utf8mb4&parseTime=true&loc=Local",
		MaxIdle:     10,
		MaxOpen:     10,
		MaxLifetime: 10,
		LogMode:     true,
	}

	db, clean := Open(dbConfig)
	defer clean()
	data := BaseApps{}
	Trace(xtrace.NewCtxWithTraceID(context.Background()), db).Table("base_apps").Where("app_id = ?", "base").First(&data)
	t.Log(data)
}
