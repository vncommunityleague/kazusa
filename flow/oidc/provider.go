package oidc

import (
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type Provider interface {
	OAuth() (*oauth2.Config, error)
	Callback(token *oauth2.Token) (string, error)
}

var providers = map[string]Provider{
	"discord": NewDiscordProvider(),
	"osu":     NewOsuProvider(),
}

func ProviderByName(name string) Provider {
	return providers[name]
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
