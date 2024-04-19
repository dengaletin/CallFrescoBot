package models

import "time"

type Campaign struct {
	Id        uint64    `gorm:"primaryKey;auto_increment"`
	Code      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
}
