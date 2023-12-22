package dao

import "gorm.io/gorm"

func MyPaging(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		switch {
		case size > 100:
			size = 100
		case size < 1:
			size = 5
		}
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
