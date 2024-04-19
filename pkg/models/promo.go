package models

import "time"

type Promo struct {
	Id         uint64    `gorm:"primaryKey;auto_increment"`
	CampaignId uint64    `gorm:"not null"`
	UserId     uint64    `gorm:"not null;unique"`
	CreatedAt  time.Time `gorm:"autoCreateTime:milli"`
}
