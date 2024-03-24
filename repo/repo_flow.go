package repo

import (
	"github.com/redis/rueidis/om"

	"github.com/vncommunityleague/kazusa/flow/oidc"
)

func (r *repositoryImpl) GetOIDCFlowRepo() om.Repository[oidc.Flow] {
	return r.OIDCFlowRepo
}
