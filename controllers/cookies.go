package controllers

import "net/http"

const (
	CookieSession = "session"
)

func NewCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
}
