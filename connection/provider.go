package connection

import (
	"context"
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/vncommunityleague/kazusa/internal"
)

type (
	providerDependencies interface {
		Repository
	}

	AuthProvider interface {
		OAuth() (*oauth2.Config, error)

		Callback(ctx context.Context, token *oauth2.Token) (*Connection, error)
	}

	SelfAuthProvider interface {
		ClientCredentials() (*clientcredentials.Config, error)
	}
)

var (
	// List of supported OAuth providers
	authProviders = map[string]func(d providerDependencies) AuthProvider{
		"osu": NewOsuAuthProvider,
	}
)

func GetAuthProvider(id string, d providerDependencies) (AuthProvider, error) {
	if p, ok := authProviders[id]; ok {
		return p(d), nil
	}

	return nil, errors.New("auth provider is not found or not supported")
}

func RouteBaseCallbackPath(provider string) string {
	return internal.GetSiteUrl() + "connections/" + provider + "/callback"
}
