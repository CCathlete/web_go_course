package controllers

import "net/http"

const (
	CookieNameSession = "session"
)

func NewCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
}

func deleteCookie(w http.ResponseWriter, name string) {
	rewrittenC := NewCookie(name, "") // Rewriting the cookie with name=name.
	rewrittenC.MaxAge = -1            // Marking the cookie for deletion at the browser.
	http.SetCookie(w, rewrittenC)
}
