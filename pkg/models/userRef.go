package models

import "time"

type UserRef struct {
	Id        uint64    `gorm:"primaryKey;auto_increment"`
	UserId1   uint64    `gorm:"not null"`
	UserId2   uint64    `gorm:"not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
}
