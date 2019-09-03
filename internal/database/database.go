package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/mkarpusiewicz/go-db-docker-boilerplate/internal/app"
)

var Connection *gorm.DB

func newConnection() (*gorm.DB, error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("database configuration failed: %v", err)
	}

	connString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", cfg.host, cfg.port, cfg.database, cfg.user, cfg.password)
	if cfg.ssl != "" {
		connString += fmt.Sprintf(" sslmode=%s", cfg.ssl)
	}

	db, err := gorm.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("database open connection failed: %v", err)
	}

	if app.Config.Env == app.DebugENV {
		db.LogMode(true)
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, fmt.Errorf("database ping failed: %v", err)
	}

	return db, nil
}

func Seed(db *gorm.DB) {
	db.
		Delete(Entry{}).
		Create(&Entry{Model: Model{ID: 1}, Name: "Test"})
}

func InitDB() error {
	db, err := newConnection()
	if err != nil {
		return fmt.Errorf("database initialization failed: %v", err)
	}

	db.AutoMigrate(&Entry{})
	Seed(db)

	Connection = db

	return nil
}
