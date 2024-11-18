package internal

import (
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

func GetUser(r *http.Request) (string, error) {
	ctx := r.Context()
	auth := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	jwksURI := os.Getenv("JWKS_URL")

	keySet := oidc.NewRemoteKeySet(ctx, jwksURI)
	verifier := oidc.NewVerifier(GetSiteUrl(), keySet, &oidc.Config{SkipClientIDCheck: true})

	token, err := verifier.Verify(ctx, auth)
	if err != nil {
		return "", err
	}

	return token.Subject, nil
}
