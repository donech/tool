//Package xdb code
package xdb

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/donech/tool/xlog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

//Open the database connection
func Open(conf Config) (*gorm.DB, func()) {
	if db, err := gorm.Open("mysql", conf.Dsn); err != nil {
		panic(errors.WithStack(err))
	} else {
		db.LogMode(false)
		db.DB().SetMaxIdleConns(conf.MaxIdle)
		db.DB().SetMaxOpenConns(conf.MaxOpen)
		db.DB().SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
		RegisterCallback(db)
		cleanup := func() {
			if err := db.Close(); err != nil {
				panic(errors.WithStack(err))
			}
		}

		return db, cleanup
	}
}

func Trace(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.InstantSet("ctx", ctx)
}

func RegisterCallback(db *gorm.DB) {
	db.Callback().Create().Register("logger", func(scope *gorm.Scope) {
		log("create", scope)
	})
	db.Callback().Update().Register("logger", func(scope *gorm.Scope) {
		log("update", scope)
	})
	db.Callback().Delete().Register("logger", func(scope *gorm.Scope) {
		log("delete", scope)
	})
	db.Callback().Query().Register("logger", func(scope *gorm.Scope) {
		log("query", scope)
	})
	db.Callback().RowQuery().Register("logger", func(scope *gorm.Scope) {
		log("raw_query", scope)
	})
}

func log(msg string, scope *gorm.Scope) {
	ctx, ok := scope.Get("ctx")
	if ok {
		xlog.L(ctx.(context.Context)).Info(
			msg,
			zap.String("sql", scope.SQL),
			zap.Reflect("value", scope.Value),
		)
	}
}
