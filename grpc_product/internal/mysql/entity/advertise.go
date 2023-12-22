package entity

import "gorm.io/gorm"

type Advertise struct {
	gorm.Model
	Index int32  `gorm:"type:int;not null;default:1"`
	Image string `gorm:"type:varchar(256);not null"`
	Url   string `gorm:"type:varchar(256);not null"`
	Sort  int32  `gorm:"type:int;not null;default:1"`
}

func (a *Advertise) TableName() string {
	return "advertise"
}
