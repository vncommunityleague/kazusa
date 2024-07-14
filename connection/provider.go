package connection

import (
	"context"
	"errors"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/vncommunityleague/kazusa/ory"
)

type (
	providerDepdencies interface {
		Repository

		ory.Provider
	}

	AuthProvider interface {
		OAuth() (*oauth2.Config, error)

		Callback(ctx context.Context, token *oauth2.Token, containers *Connections) error
	}

	SelfAuthProvider interface {
		ClientCredentials() (*clientcredentials.Config, error)
	}
)

var (
	// List of supported OAuth providers
	authProviders = map[string]func(d providerDepdencies) AuthProvider{
		"osu": NewOsuAuthProvider,
	}
)

func GetAuthProvider(id string, d providerDepdencies) (AuthProvider, error) {
	if p, ok := authProviders[id]; ok {
		return p(d), nil
	}

	return nil, errors.New("auth provider is not found or not supported")
}

func RouteBaseCallbackPath(provider string) string {
	return os.Getenv("SITE_URL") + "/connections/" + provider + "/callback"
}
