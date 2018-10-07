package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/mrl-athomelab/website/application/logger"
)

var (
	db     *gorm.DB
	tables = map[string]interface{}{
		"administrators": &Administrator{},
		"members":        &Member{},
	}
)

const (
	ByUsernamePassword = int8(iota)
	ByID
)

func Open(provider, connString string) error {
	var err error
	db, err = gorm.Open(provider, connString)
	if err != nil {
		return err
	}

	for name, object := range tables {
		logger.Info("Checking table %s ...", name)
		if !db.HasTable(object) {
			db.CreateTable(object)
		}
		db.AutoMigrate(object)
	}

	return nil
}
