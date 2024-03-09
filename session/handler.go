package session

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/vncommunityleague/kazusa/internal"
)

type (
	handlerDependencies interface {
		ManagerProvider
		Repository
	}
	HandlerProvider interface {
		SessionHandler() *Handler
	}
	Handler struct {
		r handlerDependencies
	}
)

func NewHandler(d handlerDependencies) *Handler {
	return &Handler{
		r: d,
	}
}

var (
	RouteBasePath = "/sessions"

	RouteMe      = RouteBasePath + "/me"
	RouteSession = RouteBasePath + "/:id"
)

func (h *Handler) RegisterRoutes(r *internal.Router) {
	for _, m := range []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodConnect, http.MethodOptions, http.MethodTrace} {
		r.HandleFunc(m, RouteMe, h.me)
	}

	r.DELETE(RouteSession, h.deleteMySession)
	r.DELETE(RouteBasePath, h.deleteMySessions)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	s, err := h.r.SessionManager().GetSessionFromRequest(r.Context(), r)
	if err != nil {
		internal.ErrorJson(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	internal.Json(w, http.StatusOK, s)
}

func (h *Handler) deleteMySession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "me" {
		h.me(w, r)
		return
	}

	s, err := h.r.SessionManager().GetSessionFromRequest(r.Context(), r)
	if err != nil {
		internal.ErrorJson(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	sid, err := uuid.Parse(id)
	if err != nil {
		internal.ErrorJson(w, http.StatusBadRequest, "invalid_session_id", err)
		return
	}
	if sid == s.ID {
		internal.ErrorJson(w, http.StatusBadRequest, "delete_current_session", errors.New("you cannot delete your current session"))
		return
	}

	if err := h.r.DeactivateSession(r.Context(), s.IdentityID, sid); err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, "delete_session", err)
		return
	}
}

func (h *Handler) deleteMySessions(w http.ResponseWriter, r *http.Request) {
	s, err := h.r.SessionManager().GetSessionFromRequest(r.Context(), r)
	if err != nil {
		internal.ErrorJson(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	_, err = h.r.DeactivateSessionsFromIdentityExcept(r.Context(), s.IdentityID, s.ID)
	if err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, "delete_sessions", err)
		return
	}
}
