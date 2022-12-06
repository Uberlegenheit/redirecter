package models

type Link struct {
	ID     uint64 `gorm:"column:id;PRIMARY_KEY"`
	UserID uint64 `gorm:"column:user_id"`
	Code   string `gorm:"column:code"`
}
