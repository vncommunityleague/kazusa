package oidc

import (
	"github.com/redis/rueidis/om"
)

type Repository interface {
	GetOIDCFlowRepo() om.Repository[Flow]
}
