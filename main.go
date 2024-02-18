package main

import (
	"github.com/vncommunityleague/kazusa/flow/oidc"
	"github.com/vncommunityleague/kazusa/repo"
	"github.com/vncommunityleague/kazusa/x"
	"log"
	"net/http"
	"os"
)

func main() {
	rds, err := repo.ConnectToRedis()
	if err != nil {
		panic(err)
	}

	r := repo.NewRepository(repo.RepositoryDependencies{
		Rds: rds,
	})

	router := x.NewRouter()
	oidc.NewHandler(r).RegisterRoutes(router)

	host := os.Getenv("HOST_ADDR")

	log.Println("Listening on", host)
	log.Fatal(http.ListenAndServe(host, router.ServeMux))
}
