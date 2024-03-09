package oidc

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/vncommunityleague/kazusa/internal"
	"github.com/vncommunityleague/kazusa/session"
)

type (
	Handler struct {
		d oidcDependencies
	}
)

func NewHandler(d oidcDependencies) *Handler {
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

func (h *Handler) RegisterRoutes(r *internal.Router) {
	r.GET(RouteInitPath, h.init)
	r.GET(RouteCallbackPath, h.callback)
}

func (h *Handler) init(w http.ResponseWriter, r *http.Request) {
	redirectUrl := r.FormValue("redirect")
	provider, err := GetProvider(r.PathValue("provider"), h.d)
	if err != nil {
		internal.ErrorJson(w, http.StatusNotFound, ErrOIDCProviderNotFound, err)
		return
	}

	ctx := r.Context()

	o, err := provider.OAuth()
	if err != nil {
		internal.ErrorJson(w, http.StatusNotFound, ErrOIDCProviderUnableToCreate, err)
		return
	}

	state := internal.RandomString(24)

	codeVerifier := oauth2.GenerateVerifier()
	if err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, ErrOIDCVerifierUnableToGenerate, err)
		return
	}

	if err = h.d.UpsertOIDCFlow(ctx, state, &Flow{
		CodeVerifier: codeVerifier,
		Url:          redirectUrl,
	}); err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, ErrOIDCFlowUnableToCreate, err)
		return
	}

	url := o.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(codeVerifier),
	)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	state := r.FormValue("state")
	provider, err := GetProvider(r.PathValue("provider"), h.d)
	if err != nil {
		internal.ErrorJson(w, http.StatusNotFound, ErrOIDCProviderNotFound, err)
	}

	ctx := r.Context()
	r = r.WithContext(ctx)

	flow, err := h.d.GetAndDeleteOIDCFlow(ctx, state)
	if err != nil {
		internal.ErrorJson(w, http.StatusNotFound, ErrOIDCFlowNotFound, err)
		return
	}

	t, err := exchangeCode(ctx, provider, code, flow.CodeVerifier)
	if err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, ErrOIDCUnableToExchangeCode, err)
		return
	}

	identity, created, err := provider.Callback(ctx, t)
	if err != nil {
		panic(err)
	}

	if created {
		err = requestCreateNewUser(UserCreation{
			Id:       identity.ID.String(),
			Username: "",
		})

		if err != nil {
			internal.ErrorJson(w, http.StatusInternalServerError, "unable_to_create_user", err)
			return
		}
	}

	sess, err := session.NewActiveSession(r, identity)
	if err != nil {
		panic(err)
	}

	h.d.SessionManager().IssueCookie(w, r, sess)

	if err = h.d.UpsertSession(ctx, sess); err != nil {
		panic(err)
	}

	http.Redirect(w, r, flow.Url+"?session="+sess.Token, http.StatusTemporaryRedirect)
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
