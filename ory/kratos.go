package ory

import (
	"errors"
	"os"

	"net/http"

	kratos "github.com/ory/kratos-client-go"
)

type Kratos interface {
	GetSessionFromRequest(r *http.Request) (*kratos.Session, error)
}

type (
	kratosDepdencies interface {
	}

	DefaultKratos struct {
		d kratosDepdencies
	}
)

func NewDefaultKratos(d kratosDepdencies) *DefaultKratos {
	return &DefaultKratos{
		d: d,
	}
}

const (
	Public = "public"
	Admin  = "admin"
)

func (k *DefaultKratos) getUrl(clientType string) (string, error) {
	var url string
	var exists bool

	if clientType == Admin {
		url, exists = os.LookupEnv("KRATOS_ADMIN_URL")
	} else {
		url, exists = os.LookupEnv("KRATOS_URL")
	}

	if !exists {
		if clientType == Admin {
			return "", errors.New("KRATOS_ADMIN_URL is not set")
		} else {
			return "", errors.New("KRATOS_URL is not set")
		}
	}

	return url, nil
}

func (k *DefaultKratos) getClient(clientType string) (*kratos.APIClient, error) {
	url, err := k.getUrl(clientType)
	if err != nil {
		return nil, err
	}

	cfg := kratos.NewConfiguration()
	cfg.Servers = kratos.ServerConfigurations{{URL: url}}

	return kratos.NewAPIClient(cfg), nil
}

func (k *DefaultKratos) GetSessionFromRequest(r *http.Request) (*kratos.Session, error) {
	client, err := k.getClient(Public)
	if err != nil {
		return nil, err
	}

	cookie, err := r.Cookie("ory_kratos_session")
	if err != nil {
		return nil, err
	}

	sess, _, err := client.FrontendAPI.ToSession(r.Context()).Cookie(cookie.Name + "=" + cookie.Value).Execute()
	if err != nil {
		return nil, err
	}

	return sess, nil
}
