package models

import "time"

type UserRef struct {
	UserId    uint64 `gorm:"not null;"`
	User      User
	UserRefId uint64 `gorm:"not null;unique"`
	UserRef   User
	CreatedAt time.Time `gorm:"autoCreateTime:milli;"`
}
