package game

import (
	"net/http"
	"strings"

	"github.com/vncommunityleague/kazusa/internal"
)

type (
	handlerDependenices interface {
		Repository
	}
	HandlerProvider interface {
		GameHandler() *Handler
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
	RouteBase = "/games"

	RouteGame      = RouteBase + "/{game}"
	RouteGameUsers = RouteGame + "/{users}"
)

func (h *Handler) RegisterPublicRoutes(r *internal.PublicRouter) {
	r.GET(RouteGameUsers, h.gameUsers)
}

func (h *Handler) gameUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	p, err := GetGameProvider(r.PathValue("game"), h.d)
	if err != nil {
		internal.ErrorJson(w, http.StatusNotFound, "game_provider_not_found", err)
		return
	}

	users := strings.Split(r.PathValue("users"), ",")
	mode := r.FormValue("mode")

	data, err := p.GetMultiUserGameData(ctx, users, &ExtraQuery{
		Mode: mode,
	})
	if err != nil {
		internal.ErrorJson(w, http.StatusBadRequest, "unable_to_retrieve_data", err)
		return
	}

	internal.Json(w, http.StatusOK, data)
}
