package xdb

import (
	"context"
	"database/sql"
	"testing"

	"github.com/donech/tool/xtrace"

	"github.com/donech/tool/xlog"
)

type LotteryRecord struct {
	Id          sql.NullInt64  `xorm:"pk autoincr BIGINT(20) 'id'"`
	Number      sql.NullString `xorm:"VARCHAR(255) 'number'"`
	Period      sql.NullString `xorm:"VARCHAR(255) 'period'"`
	Type        sql.NullString `xorm:"VARCHAR(255) 'type'"`
	CreatedTime sql.NullTime   `xorm:"DATETIME 'created_time'"`
	UpdatedTime sql.NullTime   `xorm:"DATETIME 'updated_time'"`
	DeletedTime sql.NullTime   `xorm:"DATETIME 'deleted_time'"`
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
		Dsn:         "root:example@tcp(localhost:3306)/nirvana?charset=utf8mb4&parseTime=true&loc=Local",
		MaxIdle:     10,
		MaxOpen:     10,
		MaxLifetime: 10,
		LogMode:     true,
	}

	db, clean := Open(dbConfig)
	defer clean()
	data := LotteryRecord{}
	Trace(xtrace.NewCtxWithTraceID(context.Background()), db).Table("lottery_record").Where("id = ?", "1").First(&data)
	t.Log(data)
}
