package connection

import (
	"context"
	"errors"
	"os"

	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/vncommunityleague/kazusa/ory"
)

type (
	Provider interface {
		OAuth() (*oauth2.Config, error)
		Callback(ctx context.Context, token *oauth2.Token) (interface{}, error)
	}

	dependencies interface {
		Repository
		ManagementProvider

		ory.Provider
	}
)

var (
	// List of supported OAuth providers
	providers = map[string]func(d dependencies) Provider{
		"osu": NewOsuProvider,
	}
)

func GetProvider(id string, d dependencies) (Provider, string, error) {
	if p, ok := providers[id]; ok {
		return p(d), id, nil
	}

	return nil, "", errors.New("connection provider is not found or not supported")
}

func RouteBaseCallbackPath(provider string) string {
	return os.Getenv("SITE_URL") + "/connections/" + provider + "/callback"
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
