package x

import "net/http"

// Copied from kratos/x/router.go
// But used http.ServeMux instead of httprouter.Router
type Router struct {
	*http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		ServeMux: http.NewServeMux(),
	}
}

func (r *Router) GET(path string, handle http.HandlerFunc) {
	r.HandleFuncWithMethod(http.MethodGet, path, handle)
}

func (r *Router) HEAD(path string, handle http.HandlerFunc) {
	r.HandleFuncWithMethod(http.MethodHead, path, handle)
}

func (r *Router) POST(path string, handle http.HandlerFunc) {
	r.HandleFuncWithMethod(http.MethodPost, path, handle)
}

func (r *Router) PUT(path string, handle http.HandlerFunc) {
	r.HandleFuncWithMethod(http.MethodPut, path, handle)
}

func (r *Router) PATCH(path string, handle http.HandlerFunc) {
	r.HandleFuncWithMethod(http.MethodPatch, path, handle)
}

func (r *Router) DELETE(path string, handle http.HandlerFunc) {
	r.HandleFuncWithMethod(http.MethodDelete, path, handle)
}

func (r *Router) HandleFuncWithMethod(method, path string, handle http.HandlerFunc) {
	r.HandleFunc(method+" "+path, handle)
}

func (r *Router) HandleFunc(path string, handle http.HandlerFunc) {
	r.ServeMux.HandleFunc(path, handle)
}

func (r *Router) Handle(path string, handler http.Handler) {
	r.ServeMux.Handle(path, handler)
}

func (r *Router) Handler(path string, handler http.Handler) {
	r.ServeMux.Handle(path, handler)
}
