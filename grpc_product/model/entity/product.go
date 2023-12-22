package entity

import "gorm.io/gorm"

type Advertise struct {
	gorm.Model
	Index string `gorm:"type:int;not null;default:1"`
	Image string `gorm:"type:varchar(256);not null"`
	Url   string `gorm:"type:varchar(256);not null"`
	Sort  int32  `gorm:"type:int;not null;default:1"`
}

func (u *Advertise) TableName() string {
	return "advertise"
}
