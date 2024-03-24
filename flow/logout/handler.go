package logout

import (
	"net/http"

	"github.com/vncommunityleague/kazusa/internal"
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
