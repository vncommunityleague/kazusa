package daemon

import (
	"context"
	"log"
	"os"

	"net/http"

	"github.com/rs/cors"

	"github.com/vncommunityleague/kazusa/internal"
	"github.com/vncommunityleague/kazusa/registry"
)

func ServePublic(r registry.Registry) {
	ctx := context.Background()

	router := internal.NewPublicRouter()
	r.RegisterPublicRoutes(ctx, router)

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
