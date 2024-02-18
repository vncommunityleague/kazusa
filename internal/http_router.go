package internal

import (
	"net/http"
	"path"
)

// Router from kratos/x/router.go
// But used http.ServeMux instead of httprouter.Router

type PublicRouter struct {
	*http.ServeMux
}

func NewPublicRouter() *PublicRouter {
	return &PublicRouter{
		ServeMux: http.NewServeMux(),
	}
}

func (r *PublicRouter) GET(path string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodGet, path, handle)
}

func (r *PublicRouter) HEAD(path string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodHead, path, handle)
}

func (r *PublicRouter) POST(path string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodPost, path, handle)
}

func (r *PublicRouter) PUT(path string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodPut, path, handle)
}

func (r *PublicRouter) PATCH(path string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodPatch, path, handle)
}

func (r *PublicRouter) DELETE(path string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodDelete, path, handle)
}

func (r *PublicRouter) HandleFunc(method, path string, handle http.HandlerFunc) {
	r.HandleFuncGeneric(method+" "+path, handle)
}

func (r *PublicRouter) HandleFuncGeneric(path string, handle http.HandlerFunc) {
	r.ServeMux.HandleFunc(path, handle)
}

func (r *PublicRouter) Handle(method, path string, handler http.Handler) {
	r.ServeMux.Handle(method+" "+path, handler)
}

func (r *PublicRouter) HandleGeneric(path string, handler http.Handler) {
	r.ServeMux.Handle(path, handler)
}

func (r *PublicRouter) Handler(path string, handler http.Handler) {
	r.ServeMux.Handle(path, handler)
}

const AdminPrefix = "/admin"

type AdminRouter struct {
	*http.ServeMux
}

func NewAdminRouter() *AdminRouter {
	return &AdminRouter{
		ServeMux: http.NewServeMux(),
	}
}

func (r *AdminRouter) GET(publicPath string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodGet, path.Join(AdminPrefix, publicPath), handle)
}

func (r *AdminRouter) HEAD(publicPath string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodHead, path.Join(AdminPrefix, publicPath), handle)
}

func (r *AdminRouter) POST(publicPath string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodPost, path.Join(AdminPrefix, publicPath), handle)
}

func (r *AdminRouter) PUT(publicPath string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodPut, path.Join(AdminPrefix, publicPath), handle)
}

func (r *AdminRouter) PATCH(publicPath string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodPatch, path.Join(AdminPrefix, publicPath), handle)
}

func (r *AdminRouter) DELETE(publicPath string, handle http.HandlerFunc) {
	r.HandleFunc(http.MethodDelete, path.Join(AdminPrefix, publicPath), handle)
}

func (r *AdminRouter) HandleFunc(method, publicPath string, handle http.HandlerFunc) {
	r.HandleFuncGeneric(method+" "+path.Join(AdminPrefix, publicPath), handle)
}

func (r *AdminRouter) HandleFuncGeneric(publicPath string, handle http.HandlerFunc) {
	r.ServeMux.HandleFunc(path.Join(AdminPrefix, publicPath), handle)
}

func (r *AdminRouter) Handle(method, publicPath string, handler http.Handler) {
	r.ServeMux.Handle(method+" "+path.Join(AdminPrefix, publicPath), handler)
}

func (r *AdminRouter) HandleGeneric(publicPath string, handler http.Handler) {
	r.ServeMux.Handle(path.Join(AdminPrefix, publicPath), handler)
}

func (r *AdminRouter) Handler(publicPath string, handler http.Handler) {
	r.ServeMux.Handle(path.Join(AdminPrefix, publicPath), handler)
}
