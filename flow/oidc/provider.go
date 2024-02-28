package oidc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/vncommunityleague/kazusa/identity"
	"github.com/vncommunityleague/kazusa/session"
)

type (
	Provider interface {
		OAuth() (*oauth2.Config, error)
		Callback(ctx context.Context, token *oauth2.Token) (*identity.Identity, error)
	}

	Dependencies interface {
		Repository

		identity.Repository

		session.Repository
		session.ManagerProvider
	}
)

var providers = map[string]func(d Dependencies) Provider{
	"discord": NewDiscordProvider,
	"osu":     NewOsuProvider,
}

func GetProvider(name string, d Dependencies) (Provider, error) {
	if p, ok := providers[name]; ok {
		return p(d), nil
	}

	return nil, errors.New("OAuth provider is not found. Please check the supported providers in the codebase")
}

func RouteBaseCallbackPath() string {
	return os.Getenv("SITE_URL") + "/flow/oidc/callback"
}

func requestOAuthUser(url string, token *oauth2.Token, output any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", "Vietnam Community League")
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return err
	}

	return nil
}
