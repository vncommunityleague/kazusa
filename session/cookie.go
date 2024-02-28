package session

import "github.com/gorilla/securecookie"

type SecureCookieProvider interface {
	SecureCookie() *securecookie.SecureCookie
}
