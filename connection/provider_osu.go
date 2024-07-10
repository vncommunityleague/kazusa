package connection

import (
	"context"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/vncommunityleague/kazusa/internal"
)

type osuAuthProvider struct {
	d providerDepdencies
}

func NewOsuAuthProvider(d providerDepdencies) AuthProvider {
	return &osuAuthProvider{
		d,
	}
}

const (
	OsuBaseURL = "https://osu.ppy.sh"
	OsuAPIURL  = OsuBaseURL + "/api/v2"

	OsuAuthURL  = OsuAPIURL + "/oauth/authorize"
	OsuTokenUrl = OsuAPIURL + "/oauth/token"

	OsuMeURL       = OsuAPIURL + "/me"
	OsuGetUsersURL = OsuAPIURL + "/users"
)

func getClientIdAndSecret() (string, string) {
	clientId := os.Getenv("OSU_CLIENT_ID")
	clientSecret := os.Getenv("OSU_CLIENT_SECRET")

	return clientId, clientSecret
}

func (p *osuAuthProvider) OAuth() (*oauth2.Config, error) {
	clientId, clientSecret := getClientIdAndSecret()

	redirectUrl := RouteBaseCallbackPath("osu")

	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  OsuAuthURL,
			TokenURL: OsuTokenUrl,
		},
		RedirectURL: redirectUrl,
		Scopes:      []string{"identify"},
	}, nil
}

func (p *osuAuthProvider) Callback(ctx context.Context, token *oauth2.Token, containers *Connections) error {
	o, err := p.OAuth()
	if err != nil {
		return err
	}

	var user struct {
		ID        uint32 `json:"id,omitempty"`
		Username  string `json:"username,omitempty"`
		AvatarURL string `json:"avatar_url,omitempty"`
		Country   string `json:"country_code,omitempty"`
	}

	httpClient := o.Client(ctx, token)
	if err := internal.RequestOAuthData(httpClient, OsuMeURL, &user); err != nil {
		return err
	}

	containers.Osu = Connection{
		Id:        user.ID,
		Username:  user.Username,
		AvatarUrl: user.AvatarURL,

		Country: user.Country,
	}

	return nil
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
		Scopes:       []string{"identify"},
	}, nil
}
