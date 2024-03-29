package models

import "time"

type User struct {
	Id           uint64 `gorm:"primaryKey;auto_increment"`
	Name         string
	TgId         int64     `gorm:"not null;index:idx_tg_id,unique"`
	IsActive     bool      `gorm:"not null;default:1"`
	CreatedAt    time.Time `gorm:"autoCreateTime:milli"`
	LastLogin    time.Time `gorm:"autoUpdateTime:milli"`
	IsNew        bool      `gorm:"not null;default:1"`
	Mode         int64     `gorm:"not null;default:0"`
	Dialog       int64     `gorm:"not null;default:0"`
	Lang         int64     `gorm:"not null;default:1"`
	DialogFromId uint64    `gorm:"not null;index:idx_dialog_from_id"`
}
