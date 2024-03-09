package repo

import (
	"github.com/vncommunityleague/kazusa/identity"
	"github.com/vncommunityleague/kazusa/session"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func ConnectToDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Debug().Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	err = db.AutoMigrate(&identity.Identity{}, &session.Session{})
	if err != nil {
		panic(err)
	}

	return db, nil
}
