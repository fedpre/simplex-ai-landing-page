package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

var (
	hashKey  = securecookie.GenerateRandomKey(64)
	blockKey = securecookie.GenerateRandomKey(32)
	s        = securecookie.New(hashKey, blockKey)
)

func GenerateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func ValidateCSRFToken(r *http.Request) bool {
	token := r.FormValue("csrf_token")
	expectedToken, err := r.Cookie("csrf_token")

	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(token), []byte(expectedToken.Value)) == 1
}

func NewSecureCookie(name, value string, expires time.Time) *http.Cookie {
	encoded, err := s.Encode(name, value)
	if err != nil {
		return nil
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	return cookie
}

func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("auth")
	if err != nil {
		return false
	}

	var value string
	if err := s.Decode("auth", cookie.Value, &value); err != nil {
		return false
	}

	return value == "authenticated"
}
