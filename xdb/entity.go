package xdb

type Entity struct {
	ID          int64 `gorm:"primary_key" json:"id"`
	CreatedTime int   `json:"created_time"`
	UpdatedTime int   `json:"updated_time"`
	DeletedTime int   `json:"deleted_time"`
}
