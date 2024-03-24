package repo

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"
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

		OIDCFlowRepo om.Repository[oidc.Flow]
	}
)

func NewRepository(d RepositoryDependencies) Repository {
	oidcFlowRepo := om.NewJSONRepository[oidc.Flow]("oidc_flow", oidc.Flow{}, d.Rds)

	return &repositoryImpl{
		d:            d,
		OIDCFlowRepo: oidcFlowRepo,
	}
}

func (r *repositoryImpl) Raw(query string, args ...interface{}) {
	r.d.DB.Raw(query, args)
}

func (r *repositoryImpl) Exec(query string, args ...interface{}) {
	r.d.DB.Exec(query, args)
}
