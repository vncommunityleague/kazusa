package repo

import (
	"os"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/vncommunityleague/kazusa/connection"
)

type (
	Dependencies struct {
		Rds rueidis.Client
		DB  *gorm.DB
	}

	Repository interface {
		connection.Repository

		Raw(query string, args ...interface{})
		Exec(query string, args ...interface{})
	}

	Default struct {
		d Dependencies

		ConnectionFlowRepo om.Repository[connection.Flow]
	}
)

func connectToDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Debug().Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	err = db.AutoMigrate(&connection.Connection{})
	if err != nil {
		panic(err)
	}

	return db, nil
}

func connectToRedis() (rueidis.Client, error) {
	url := os.Getenv("REDIS_URL")

	client, err := rueidis.NewClient(rueidis.MustParseURL(url))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewRepository() Repository {
	db, err := connectToDB()
	if err != nil {
		panic(err)
	}

	rds, err := connectToRedis()
	if err != nil {
		panic(err)
	}

	return newRepository(Dependencies{
		Rds: rds,
		DB:  db,
	})
}

func newRepository(d Dependencies) Repository {
	connectionFlowRepo := om.NewJSONRepository[connection.Flow]("connection_flow", connection.Flow{}, d.Rds)

	return &Default{
		d:                  d,
		ConnectionFlowRepo: connectionFlowRepo,
	}
}

func (r *Default) Raw(query string, args ...interface{}) {
	r.d.DB.Raw(query, args)
}

func (r *Default) Exec(query string, args ...interface{}) {
	r.d.DB.Exec(query, args)
}
