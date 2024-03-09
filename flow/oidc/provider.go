package oidc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"net/http"
	"os"

	"github.com/vncommunityleague/kazusa/identity"
	"github.com/vncommunityleague/kazusa/session"
)

type (
	Provider interface {
		OAuth() (*oauth2.Config, error)
		Callback(ctx context.Context, token *oauth2.Token) (*identity.Identity, bool, error)
	}

	oidcDependencies interface {
		Repository

		identity.Repository

		session.Repository
		session.ManagerProvider
	}
)

var (
	// List of supported OAuth providers
	providers = map[string]func(d oidcDependencies) Provider{
		"discord": NewDiscordProvider,
		"osu":     NewOsuProvider,
	}
)

func GetProvider(id string, d oidcDependencies) (Provider, error) {
	if p, ok := providers[id]; ok {
		return p(d), nil
	}

	return nil, errors.New("OIDC provider is not found or not supported")
}

func RouteBaseCallbackPath() string {
	return os.Getenv("SITE_URL") + "/flow/oidc/callback"
}

type UserCreation struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func requestCreateNewUser(data UserCreation) error {
	url := os.Getenv("USER_SERVICE_URL") + "/internal/_create_new_user"
	token := os.Getenv("USER_SERVICE_TOKEN")

	userCreationJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(userCreationJson))
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "Vietnam Community League")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unable to create new user")
	}

	return nil
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
