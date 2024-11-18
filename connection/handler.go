package connection

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/vncommunityleague/kazusa/internal"
)

type (
	handlerDependencies interface {
		Repository
	}
	HandlerProvider interface {
		ConnectionHandler() *Handler
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

const ExpiredTime = 5 * time.Minute

type Flow struct {
	Key       string    `json:"key" redis:",key"`
	Ver       int64     `json:"ver" redis:",ver"`
	ExpiresAt time.Time `json:"exat" redis:",exat"`

	UserId       string `json:"user_id"`
	CodeVerifier string `json:"code_verifier"`
	Url          string `json:"url"`
}

var (
	RouteBase = "/connections"

	RouteAuthorize = RouteBase + "/{provider}/authorize"
	RouteCallback  = RouteBase + "/{provider}/callback"
)

func (h *Handler) RegisterPublicRoutes(r *internal.PublicRouter) {
	r.GET(RouteAuthorize, h.authorize)
	r.GET(RouteCallback, h.callback)
}

//func (h *Handler) RegisterAdminRoutes(r *internal.AdminRouter) {}

func (h *Handler) authorize(w http.ResponseWriter, r *http.Request) {
	redirectUrl := r.FormValue("redirect")
	provider, err := GetAuthProvider(r.PathValue("provider"), h.d)
	if err != nil {
		internal.ErrorJson(w, http.StatusNotFound, "connection_provider_not_found", err)
		return
	}
	o, err := provider.OAuth()
	if err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, "connection_provider_not_found", err)
		return
	}

	ctx := r.Context()

	userId, err := internal.GetUser(r)
	if err != nil {
		internal.ErrorJson(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	state := internal.RandomString(24)
	codeVerifier := oauth2.GenerateVerifier()

	if err = h.d.GetConnectionFlows().Save(ctx, &Flow{
		Key:          state,
		UserId:       userId,
		CodeVerifier: codeVerifier,
		Url:          redirectUrl,
		ExpiresAt:    time.Now().Add(ExpiredTime),
	}); err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, "unable_to_init_flow", err)
		return
	}

	url := o.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(codeVerifier),
	)

	println(url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	state := r.FormValue("state")
	provider, err := GetAuthProvider(r.PathValue("provider"), h.d)
	if err != nil {
		internal.ErrorJson(w, http.StatusNotFound, "connection_provider_not_found", err)
	}

	ctx := r.Context()

	userId, err := internal.GetUser(r)
	if err != nil {
		internal.ErrorJson(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	flow, err := h.d.GetConnectionFlows().Fetch(ctx, state)
	if err != nil {
		internal.ErrorJson(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	if flow.UserId != userId {
		internal.ErrorJson(w, http.StatusUnauthorized, "user_not_match", errors.New("user is not matched"))
		return
	}

	id, err := uuid.Parse(userId)
	if err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, "uuid_cannot_be_converted", err)
		return
	}

	t, err := exchangeCode(ctx, provider, code, flow.CodeVerifier)
	if err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, "unable_to_exchange_code", err)
		return
	}

	con, err := provider.Callback(ctx, t)
	if err != nil {
		internal.ErrorJson(w, http.StatusBadRequest, "unable_to_link", err)
		return
	}

	con.UserId = id

	if err := h.d.SaveConnection(ctx, con); err != nil {
		internal.ErrorJson(w, http.StatusBadRequest, "unable_to_link", err)
		return
	}

	http.Redirect(w, r, flow.Url, http.StatusTemporaryRedirect)
}

func exchangeCode(ctx context.Context, provider AuthProvider, code string, codeVerifier string) (*oauth2.Token, error) {
	o, err := provider.OAuth()
	if err != nil {
		return nil, err
	}

	t, err := o.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
	if err != nil {
		return nil, err
	}

	return t, nil
}
