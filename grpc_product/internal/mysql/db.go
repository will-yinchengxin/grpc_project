package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"grpc_product/internal/config"
	"grpc_product/internal/mysql/entity"
	"log"
	"os"
	"time"
)

var (
	MysqlPool map[string]*gorm.DB
)

func InitMysql() {
	MysqlPool = make(map[string]*gorm.DB)
	db, err := initGorm(config.AppConf.DBConfig)
	if err != nil {
		panic("init db err " + err.Error())
	}
	if err == nil && db != nil {
		MysqlPool[config.AppConf.DBConfig.DBName] = db
	}
}

func initGorm(cfg config.DBConfig) (*gorm.DB, error) {
	connectStr := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connectStr,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second * 3,
				LogLevel:      logger.Silent,
				Colorful:      false,
			},
		),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	//db = db.Debug() // start debug mod

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetConnMaxLifetime(time.Second * 28800)
	sqlDB.SetConnMaxIdleTime(time.Second * 7200)
	dbMigrate(db)
	log.Println("init mysql success")
	return db, nil
}

func dbMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&entity.Advertise{},
	)
}

func GetDB() (db *gorm.DB, err error) {
	return MysqlPool[config.AppConf.DBConfig.DBName], nil
}
