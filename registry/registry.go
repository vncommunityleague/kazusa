package registry

import (
	"github.com/vncommunityleague/kazusa/repo"
	"github.com/vncommunityleague/kazusa/session"
)

type Registry interface {
	repo.Repository

	session.SecureCookieProvider
	session.ManagerProvider
}
