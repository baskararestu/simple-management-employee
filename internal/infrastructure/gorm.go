package infrastructure

import (
	"simple-management-employee/internal/domain"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func dbSetup() {
	var err error
	l := gormLogger.Default.LogMode(gormLogger.Silent)
	if cfg.Database.Driver == "mysql" {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN: cfg.Database.DSN,
		}), &gorm.Config{
			Logger: l,
		})
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{
			Logger: l,
		})
	}

	if err != nil {
		panic(err)
	}

	if cfg.IsDevelopment {
		if err := db.AutoMigrate(
			&domain.User{},
			&domain.Role{},
			&domain.AnnualLeave{},
		); err != nil {
			panic(err)
		}
	}
}
