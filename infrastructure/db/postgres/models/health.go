package models

import "gorm.io/gorm"

const TableNamePing = "pings"

type Ping struct {
	gorm.Model
	ID string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`

	Name string `gorm:"not null"`
}

func (p *Ping) TableName() string {
	return TableNamePing
}
