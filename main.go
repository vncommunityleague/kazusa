package main

import (
	"github.com/rs/cors"
	"github.com/vncommunityleague/kazusa/registry"
	"github.com/vncommunityleague/kazusa/session"
	"log"
	"net/http"
	"os"

	"github.com/vncommunityleague/kazusa/flow/oidc"
	"github.com/vncommunityleague/kazusa/internal"
)

func main() {
	reg := registry.NewRegistryDefault()
	router := internal.NewRouter()
	oidc.NewHandler(reg).RegisterRoutes(router)
	session.NewHandler(reg).RegisterRoutes(router)

	host := os.Getenv("HOST_ADDR")

	cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Cookie", "Content-Type"},
		ExposedHeaders:   []string{"Content-Type", "Set-Cookie"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler(router)

	log.Println("Listening on", host)
	log.Fatal(http.ListenAndServe(host, router.ServeMux))
}
