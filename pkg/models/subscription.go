package models

import (
	"encoding/json"
	"time"
)

type Subscription struct {
	Id          uint64 `gorm:"primaryKey;auto_increment"`
	UserId      uint64 `gorm:"not null;index:idx_user_id__active_due"`
	User        User
	PlanId      *uint64 `gorm:"default:null"`
	Plan        Plan
	Limit       int       `gorm:"not null;"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli"`
	ActiveDue   time.Time `gorm:"index:idx_user_id__active_due"`
	RefreshDate time.Time
	Usage       json.RawMessage `gorm:"type:json;default:null"`
}
