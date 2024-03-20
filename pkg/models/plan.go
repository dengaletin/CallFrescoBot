package models

import "encoding/json"

type Plan struct {
	Id     uint64          `gorm:"primaryKey;auto_increment"`
	Name   string          `gorm:"not null"`
	Config json.RawMessage `gorm:"type:json"`
}
