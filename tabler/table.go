package tabler

import (
	"github.com/jinzhu/gorm"
)

type table struct {
	db        *gorm.DB
	cursorKey string
	Pager
}

type Pager struct {
	PageSize int32
	PageNum  int32
	Cursor   int64
}

func NewTable(db *gorm.DB, cursorKey string, pager Pager) *table {
	return &table{
		db:        db,
		cursorKey: cursorKey,
		Pager:     pager,
	}
}

func (t table) Build() (res *gorm.DB) {
	res = t.db
	if t.PageNum == 0 {
		res = t.db.Where("? < ?", t.cursorKey, t.Cursor)
		res = res.Limit(t.PageSize)
		return res
	}
	offset := t.PageSize * (t.PageNum - 1)
	return res.Offset(offset).Limit(t.PageSize)
}

func EmptyDBFuc(db *gorm.DB) *gorm.DB {
	return db
}

func (t table) FindAll(doer func(db *gorm.DB) *gorm.DB) error {
	return doer(t.Build()).Error
}
