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
	return &repositoryImpl{
		d,
	}
}

func (r *repositoryImpl) Raw(query string, args ...interface{}) {
	r.d.DB.Raw(query, args)
}

func (r *repositoryImpl) Exec(query string, args ...interface{}) {
	r.d.DB.Exec(query, args)
}
