package logout

import (
	"github.com/vncommunityleague/kazusa/internal"
	"net/http"
)

type Handler struct {
}

func (h *Handler) RegisterRoutes(router internal.Router) {

}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	sess, err := h.d.SessionManager().FetchFromRequest(r.Context(), r)
	if err != nil {
		h.d.SelfServiceErrorManager().Forward(r.Context(), w, r, err)
		return
	}
}
