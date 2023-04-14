package models

import "time"

type Message struct {
	Id        uint64 `gorm:"primaryKey;auto_increment"`
	UserId    uint64 `gorm:"not null;index:idx_user_id__created_at"`
	User      User
	Message   string    `gorm:"not null"`
	Response  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli;index:idx_user_id__created_at"`
}
