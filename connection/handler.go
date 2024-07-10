package connection

import (
	"net/http"

	"github.com/vncommunityleague/kazusa/internal"
	"github.com/vncommunityleague/kazusa/ory"
)

type (
	handlerDependenices interface {
		Repository
		ManagementProvider

		ory.Provider
	}
	HandlerProvider interface {
		ConnectionHandler() *Handler
	}
	Handler struct {
		d handlerDependenices
	}
)

func NewHandler(d handlerDependenices) *Handler {
	return &Handler{
		d,
	}
}

var (
	RouteBaseConnections = "/connections"

	RouteConnectionsMe = RouteBaseConnections + "/me"

	RouteSettingCallback = RouteBaseConnections + "/setting-callback"

	RouteAuthorize = RouteBaseConnections + "/{provider}/authorize"
	RouteCallback  = RouteBaseConnections + "/{provider}/callback"
)

func (h *Handler) RegisterPublicRoutes(r *internal.PublicRouter) {
	r.GET(RouteConnectionsMe, h.connectionsMe)

	r.GET(RouteAuthorize, h.authorize)
	r.GET(RouteCallback, h.callback)
}

func (h *Handler) RegisterAdminRoutes(r *internal.AdminRouter) {
	r.POST(RouteSettingCallback, h.settingCallback)
}

func (h *Handler) connectionsMe(w http.ResponseWriter, r *http.Request) {
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
