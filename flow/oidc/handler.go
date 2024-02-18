package oidc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net/http"

	"github.com/vncommunityleague/kazusa/x"
)

type (
	handlerDependencies interface {
		Repository
	}

	Handler struct {
		d handlerDependencies
	}
)

func NewHandler(d handlerDependencies) *Handler {
	return &Handler{
		d,
	}
}

type Flow struct {
	CodeVerifier string `json:"code_verifier"`
	Url          string `json:"url"`
}

var (
	RouteBasePath = "/flow/oidc"

	RouteInitPath     = RouteBasePath + "/init/{provider}"
	RouteCallbackPath = RouteBasePath + "/callback/{provider}"
)

func (h *Handler) RegisterRoutes(r *x.Router) {
	r.GET(RouteInitPath, h.init)
	r.GET(RouteCallbackPath, h.callback)
}

func (h *Handler) init(w http.ResponseWriter, r *http.Request) {
	provider := ProviderByName(r.PathValue("provider"))
	redirectUrl := r.FormValue("redirect")

	ctx := r.Context()

	o, err := provider.OAuth()
	if err != nil {
		panic(err)
	}

	state, err := randomBytesInHex(24)
	if err != nil {
		panic(err)
	}

	codeVerifier := oauth2.GenerateVerifier()
	if err != nil {
		panic(err)
	}

	if err = h.d.UpsertFlow(ctx, state, &Flow{
		CodeVerifier: codeVerifier,
		Url:          redirectUrl,
	}); err != nil {
		panic(err)
	}

	url := o.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(codeVerifier),
	)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func randomBytesInHex(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic("Impl error handler for randomBytesInHex")
	}

	return hex.EncodeToString(buf), nil
}

func (h *Handler) callback(w http.ResponseWriter, r *http.Request) {
	provider := ProviderByName(r.PathValue("provider"))
	code := r.FormValue("code")
	state := r.FormValue("state")

	ctx := r.Context()
	r = r.WithContext(ctx)

	flow, err := h.d.GetFlow(ctx, state)
	if err != nil {
		panic(err)
	}

	t, err := exchangeCode(ctx, provider, code, flow.CodeVerifier)
	if err != nil {
		panic("Impl error handler for exchangeCode")
	}

	userId, err := provider.Callback(t)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "ID: %s", userId)
}

func exchangeCode(ctx context.Context, provider Provider, code string, codeVerifier string) (*oauth2.Token, error) {
	o, err := provider.OAuth()
	if err != nil {
		panic(err)
	}

	t, err := o.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
	if err != nil {
		panic(err)
	}

	return t, nil
}
