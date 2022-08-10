package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/zipzap11/pharm-API/config"
	"github.com/zipzap11/pharm-API/model"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", config.GetDBUser(), config.GetDBPassword(), config.GetDBHost(), config.GetDBName())
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
	// 	config.GetDBHost(),
	// 	config.GetDBUser(),
	// 	config.GetDBPassword(),
	// 	config.GetDBName(),
	// 	config.GetDBPort(),
	// )
	log.Info("dsn = ", dsn)
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		logrus.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}))
	if err != nil {
		logrus.Fatal(err)
	}
	AutoMigrate()
	logrus.Info("success create DB connection to postgres")
}

func AutoMigrate() {
	DB.AutoMigrate(model.ModelList()...)
}
