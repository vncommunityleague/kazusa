package repo

import (
	"github.com/redis/rueidis"
	"github.com/vncommunityleague/kazusa/flow/oidc"
	"github.com/vncommunityleague/kazusa/identity"
	"github.com/vncommunityleague/kazusa/session"
	"gorm.io/gorm"
)

type (
	RepositoryDependencies struct {
		Rds rueidis.Client
		DB  *gorm.DB
	}

	Repository interface {
		oidc.Repository
		identity.Repository
		session.Repository

		Raw(query string, args ...interface{})
		Exec(query string, args ...interface{})
	}

	repositoryImpl struct {
		d RepositoryDependencies
	}
)

func NewRepository(d RepositoryDependencies) Repository {
	r := &repositoryImpl{
		d,
	}

	r.d.DB.Debug().Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	err := r.d.DB.AutoMigrate(&identity.Identity{}, &session.Session{})
	if err != nil {
		panic(err)
	}

	return r
}

func (r *repositoryImpl) Raw(query string, args ...interface{}) {
	r.d.DB.Raw(query, args)
}

func (r *repositoryImpl) Exec(query string, args ...interface{}) {
	r.d.DB.Exec(query, args)
}
