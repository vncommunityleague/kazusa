package connection

import (
	"net/http"

	"github.com/vncommunityleague/kazusa/internal"
)

type (
	HandlerProvider interface {
		ConnectionHandler() *Handler
	}
	Handler struct {
		d dependencies
	}
)

func NewHandler(d dependencies) *Handler {
	return &Handler{
		d,
	}
}

var (
	RouteBasePath = "/connections"

	RouteMe = RouteBasePath + "/me"

	RouteAuthorizePath = RouteBasePath + "/{provider}/authorize"
	RouteCallbackPath  = RouteBasePath + "/{provider}/callback"

	RouteSettingCallback = RouteBasePath + "/setting-callback"
)

func (h *Handler) RegisterPublicRoutes(r *internal.PublicRouter) {
	r.GET(RouteMe, h.me)

	r.GET(RouteAuthorizePath, h.authorize)
	r.GET(RouteCallbackPath, h.callback)
}

func (h *Handler) RegisterAdminRoutes(r *internal.AdminRouter) {
	r.POST(RouteSettingCallback, h.settingCallback)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sess, err := h.d.Kratos().GetSessionFromRequest(r)
	if err != nil {
		internal.ErrorJson(w, http.StatusBadRequest, "session_error", err)
		return
	}

	conns, err := h.d.GetConnectionsByID(ctx, sess.Identity.Id)
	if err != nil {
		internal.ErrorJson(w, http.StatusInternalServerError, "profile_error", err)
		return
	}

	internal.Json(w, http.StatusOK, conns)
}

func (h *Handler) settingCallback(w http.ResponseWriter, r *http.Request) {
}
