package registry

import (
	"github.com/gorilla/securecookie"
	"github.com/vncommunityleague/kazusa/repo"
	"github.com/vncommunityleague/kazusa/session"
	"os"
)

type Default struct {
	repo.Repository

	sessionM *session.Manager

	sCookie *securecookie.SecureCookie
}

func NewRegistryDefault() Registry {
	rds, err := repo.ConnectToRedis()
	if err != nil {
		panic(err)
	}

	db, err := repo.ConnectToDB()
	if err != nil {
		panic(err)
	}

	return &Default{
		Repository: repo.NewRepository(repo.RepositoryDependencies{
			Rds: rds,
			DB:  db,
		}),
	}
}

func (r *Default) SecureCookie() *securecookie.SecureCookie {
	if r.sCookie == nil {
		secretKey := os.Getenv("SECRET_COOKIE")
		r.sCookie = securecookie.New([]byte(secretKey), nil)
	}

	return r.sCookie
}

func (r *Default) SessionManager() *session.Manager {
	if r.sessionM == nil {
		r.sessionM = session.NewManager(r)
	}

	return r.sessionM
}
