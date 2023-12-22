package dao

import "gorm.io/gorm"

func PageList(pageNo, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNo < 1 {
			pageNo = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize < 1:
			pageSize = 5
		}
		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
