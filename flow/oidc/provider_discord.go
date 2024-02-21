package oidc

import (
	"context"
	"github.com/vncommunityleague/kazusa/identity"
	"os"

	"golang.org/x/oauth2"
)

type DiscordProvider struct {
	d Dependencies
}

func NewDiscordProvider(d Dependencies) Provider {
	return &DiscordProvider{
		d,
	}
}

const (
	discordBaseUrl = "https://discord.com"
	discordApiUrl  = discordBaseUrl + "/api"
)

var discordEndpoint = oauth2.Endpoint{
	AuthURL:  discordBaseUrl + "/oauth2/authorize",
	TokenURL: discordApiUrl + "/oauth2/token",
}

func (p *DiscordProvider) OAuth() (*oauth2.Config, error) {
	clientId := os.Getenv("DISCORD_CLIENT_ID")
	clientSecret := os.Getenv("DISCORD_CLIENT_SECRET")

	redirectUrl := RouteBaseCallbackPath() + "/discord"

	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     discordEndpoint,
		RedirectURL:  redirectUrl,
		Scopes:       []string{"identify"},
	}, nil
}

func (p *DiscordProvider) Callback(ctx context.Context, token *oauth2.Token) (*identity.Identity, error) {
	var user struct {
		ID       string `json:"id,omitempty"`
		Username string `json:"username,omitempty"`
		Avatar   string `json:"avatar,omitempty"`
	}

	if err := requestOAuthUser(discordApiUrl+"/users/@me", token, &user); err != nil {
		return nil, err
	}

	id, err := p.d.GetIdentityByDiscordID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return id, nil
}
