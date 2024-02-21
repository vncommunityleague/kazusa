package main

import (
	"github.com/vncommunityleague/kazusa/flow/oidc"
	"github.com/vncommunityleague/kazusa/internal"
	"github.com/vncommunityleague/kazusa/repo"
	"log"
	"net/http"
	"os"
)

func main() {
	rds, err := repo.ConnectToRedis()
	if err != nil {
		panic(err)
	}

	db, err := repo.ConnectToDB()
	if err != nil {
		panic(err)
	}

	r := repo.NewRepository(repo.RepositoryDependencies{
		Rds: rds,
		DB:  db,
	})

	router := internal.NewRouter()
	oidc.NewHandler(r).RegisterRoutes(router)

	host := os.Getenv("HOST_ADDR")

	log.Println("Listening on", host)
	log.Fatal(http.ListenAndServe(host, router.ServeMux))
}
