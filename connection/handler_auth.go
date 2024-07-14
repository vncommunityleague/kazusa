package connection

import (
	"context"
	"errors"
	"time"

	"net/http"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/vncommunityleague/kazusa/internal"
)

const ExpiredTime = 5 * time.Minute

func (h *Handler) authorize(w http.ResponseWriter, r *http.Request) {
	redirectUrl := r.FormValue("redirect")
	provider, err := GetAuthProvider(r.PathValue("provider"), h.d)
	if err != nil {
		internal.ErrorJson(w, http.StatusNotFound, "connection_provider_not_found", err)
		return
	}

	ctx := r.Context()

	o, err := provider.OAuth()
	if err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, "connection_provider_not_found", err)
		return
	}

	sess, err := h.d.Kratos().GetSessionFromRequest(r)
	if err != nil {
		internal.ErrorJson(w, http.StatusUnauthorized, "session_not_found", err)
		return
	}

	state := internal.RandomString(24)
	codeVerifier := oauth2.GenerateVerifier()

	if err = h.d.GetConnectionFlows().Save(ctx, &Flow{
		Key:          state,
		SessionId:    sess.Id,
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

	sess, err := h.d.Kratos().GetSessionFromRequest(r)
	if err != nil {
		internal.ErrorJson(w, http.StatusUnauthorized, "session_not_found", err)
		return
	}

	flow, err := h.d.GetConnectionFlows().Fetch(ctx, state)
	if err != nil {
		internal.ErrorJson(w, http.StatusNotFound, "connection_flow_not_found", err)
		return
	}

	if flow.SessionId != sess.Id {
		internal.ErrorJson(w, http.StatusUnauthorized, "session_not_match", errors.New("session is not matched"))
		return
	}

	id, err := uuid.Parse(sess.Identity.Id)
	if err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, "uuid_cannot_be_converted", err)
		return
	}

	t, err := exchangeCode(ctx, provider, code, flow.CodeVerifier)
	if err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, "unable_to_exchange_code", err)
		return
	}

	containers := &Connections{
		ID: id,
	}
	if err := provider.Callback(ctx, t, containers); err != nil {
		internal.ErrorJson(w, http.StatusBadRequest, "unable_to_link", err)
		return
	}

	if err := h.d.SaveConnections(ctx, containers); err != nil {
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
