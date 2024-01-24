package models

import "time"

type Invoice struct {
	Id              uint64 `gorm:"primaryKey;auto_increment"`
	PaymentMethodId uint64 `gorm:"not null;index:idx_payment_method_id"`
	UserId          uint64 `gorm:"not null;index:idx_user_id__created_at"`
	User            User
	Amount          float64 `gorm:"not null;index:idx_amount"`
	Currency        string
	Coin            int
	Status          int64     `gorm:"not null;default:0"`
	CreatedAt       time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime:milli"`
}
