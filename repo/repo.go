package repo

import (
	"github.com/redis/rueidis"
	"github.com/vncommunityleague/kazusa/flow/oidc"
)

type (
	RepositoryDependencies struct {
		Rds rueidis.Client
	}

	Repository interface {
		oidc.Repository
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
