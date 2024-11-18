package connection

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/vncommunityleague/kazusa/internal"
)

type osuAuthProvider struct {
	d providerDependencies
}

func NewOsuAuthProvider(d providerDependencies) AuthProvider {
	return &osuAuthProvider{
		d,
	}
}

const (
	OsuBaseURL = "https://osu.ppy.sh"
	OsuAPIURL  = OsuBaseURL + "/api/v2"

	OsuAuthURL  = OsuBaseURL + "/oauth/authorize"
	OsuTokenUrl = OsuBaseURL + "/oauth/token"

	OsuMeURL = OsuAPIURL + "/me"
)

func getClientIdAndSecret() (string, string) {
	clientId := os.Getenv("OSU_CLIENT_ID")
	clientSecret := os.Getenv("OSU_CLIENT_SECRET")

	return clientId, clientSecret
}

func (p *osuAuthProvider) OAuth() (*oauth2.Config, error) {
	clientId, clientSecret := getClientIdAndSecret()

	redirectUrl := RouteBaseCallbackPath("osu")

	println(redirectUrl)

	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  OsuAuthURL,
			TokenURL: OsuTokenUrl,
		},
		RedirectURL: redirectUrl,
		Scopes:      []string{"public", "identify"},
	}, nil
}

func (p *osuAuthProvider) Callback(ctx context.Context, token *oauth2.Token) (*Connection, error) {
	o, err := p.OAuth()
	if err != nil {
		return nil, err
	}

	var user struct {
		ID        uint32 `json:"id,omitempty"`
		Username  string `json:"username,omitempty"`
		AvatarURL string `json:"avatar_url,omitempty"`
	}

	req, err := http.NewRequest(http.MethodGet, OsuMeURL, nil)
	if err != nil {
		return nil, err
	}

	httpClient := o.Client(ctx, token)
	if err := internal.RequestOAuthData(httpClient, req, &user); err != nil {
		return nil, err
	}

	return &Connection{
		ConnId:    strconv.FormatInt(int64(user.ID), 10),
		Username:  user.Username,
		AvatarUrl: user.AvatarURL,
	}, nil
}

type osuSelfAuthProvider struct{}

func NewOsuSelfAuthProvider() SelfAuthProvider {
	return &osuSelfAuthProvider{}
}

func (p *osuSelfAuthProvider) ClientCredentials() (*clientcredentials.Config, error) {
	clientId, clientSecret := getClientIdAndSecret()

	return &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     OsuTokenUrl,
		Scopes:       []string{"public"},
	}, nil
}
