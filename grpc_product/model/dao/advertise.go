package dao

import (
	"context"
	"gorm.io/gorm"
	"grpc_product/internal/mysql"
	"grpc_product/internal/mysql/entity"
	"log"
)

type Advertise struct {
	DB *gorm.DB
}

func AdvertiseWithContext(ctx context.Context) (*Advertise, error) {
	db, err := mysql.GetDB()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &Advertise{DB: db}, nil
}

func (a *Advertise) List() ([]entity.Advertise, int64, error) {
	var (
		list  = []entity.Advertise{}
		count int64
	)
	err := a.DB.Model(entity.Advertise{}).
		Count(&count).
		Order("id DESC").
		Find(&list).Error
	return list, count, err
}
