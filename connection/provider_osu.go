package connection

import (
	"context"
	"os"

	"golang.org/x/oauth2"
)

type OsuProvider struct {
	d dependencies
}

func NewOsuProvider(d dependencies) Provider {
	return &OsuProvider{
		d,
	}
}

const (
	osuBaseUrl = "https://osu.ppy.sh"
	osuApiUrl  = osuBaseUrl + "/api/v2"
)

var osuEndpoint = oauth2.Endpoint{
	AuthURL:  osuBaseUrl + "/oauth/authorize",
	TokenURL: osuBaseUrl + "/oauth/token",
}

func (p *OsuProvider) OAuth() (*oauth2.Config, error) {
	clientId := os.Getenv("OSU_CLIENT_ID")
	clientSecret := os.Getenv("OSU_CLIENT_SECRET")

	redirectUrl := RouteBaseCallbackPath("osu")

	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     osuEndpoint,
		RedirectURL:  redirectUrl,
		Scopes:       []string{"identify"},
	}, nil
}

func (p *OsuProvider) Callback(ctx context.Context, token *oauth2.Token) (interface{}, error) {
	var user struct {
		ID        uint32 `json:"id,omitempty"`
		Username  string `json:"username,omitempty"`
		AvatarURL string `json:"avatar_url,omitempty"`
	}

	if err := requestOAuthUser(osuApiUrl+"/me", token, &user); err != nil {
		return nil, err
	}

	return OsuConnection{
		Id: user.ID,
		Username: user.Username,
		AvatarUrl: user.AvatarURL,
	}, nil
}
