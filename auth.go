package auth

import (
	"io"
	"net/http"
	"strings"
)

// Matcher represents the required auth header matcher function
// signature implemented by matchers.
type Matcher func([]Token, Token) bool

// BasicAuth represents the user-password pair used as helper struct.
type BasicAuth struct {
	User, Password string
}

// Token represents the parsed schema of an HTTP authorization header.
type Token struct {
	// Type stores the authorization type. Usually Basic or Bearer.
	Type string

	// Value stores the authorization token.
	// If the authorization token has no type, this field will be filled instead.
	Value string
}

// Config represents the authorization middleware configuration.
type Config struct {
	// RealM token used for unauthorized
	RealM string

	// Tokens stores a list of allowed authorization tokens.
	// Tokens could be of any type.
	Tokens []Token

	// Matchers stores a list of matchers used to compare an incoming
	// authorization header againts the registered allowed tokens.
	Matchers []Matcher
}

// Handler authentication middleware handler.
type Handler struct {
	cfg *Config
}

// HandleHTTP will be called each time the request hits the location with this middleware activated
func (a *Handler) HandleHTTP(w http.ResponseWriter, r *http.Request, h http.Handler) {
	// If no tokens registered, just ignore the auth process
	if len(a.cfg.Tokens) == 0 {
		h.ServeHTTP(w, r)
		return
	}

	// Parse and decode auth header
	token, err := ParseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		unauthorized(w, token, a.cfg.RealM)
		return
	}

	// Compare token againsts matchers
	if !a.match(token) {
		unauthorized(w, token, a.cfg.RealM)
		return
	}

	// If authorized, pass the request to the next middleware in chain
	h.ServeHTTP(w, r)
}

// match matches a given token againts the configured matcher functions.
func (a *Handler) match(token Token) bool {
	for _, matcher := range a.cfg.Matchers {
		if matcher(a.cfg.Tokens, token) {
			return true
		}
	}
	return false
}

// New creates a new authorization middleware handler with the given config.
func New(cfg *Config) *Handler {
	// Define defaults
	if len(cfg.Matchers) == 0 {
		cfg.Matchers = []Matcher{MatchAuthHeader}
	}
	if cfg.RealM == "" {
		cfg.RealM = "Authorization required"
	}
	return &Handler{cfg}
}

// User function is optional but handy, used to check input parameters when creating new middlewares
func User(user, password string) *Handler {
	return Users(BasicAuth{User: user, Password: password})
}

// Users creates a new auth middleware allowing access to the list of users.
func Users(users ...BasicAuth) *Handler {
	tokens := make([]Token, len(users))
	for _, user := range users {
		tokens = append(tokens, Token{Type: "Basic", Value: user.User + ":" + user.Password})
	}
	return New(&Config{Tokens: tokens})
}

// Tokens creates a new auth handler allowing the given token strings.
// Tokens must be transported via Authorization header.
func Tokens(tokens ...string) *Handler {
	store := make([]Token, len(tokens))
	for _, token := range tokens {
		store = append(store, Token{Value: token})
	}
	return New(&Config{Tokens: store})
}

// MatchAuthHeader matches an authorization header againts the allowed tokens.
func MatchAuthHeader(tokens []Token, token Token) bool {
	for _, t := range tokens {
		if strings.ToLower(t.Type) == token.Type && t.Value == token.Value {
			return true
		}
	}
	return false
}

// unauthorized is used to write a forbidden status response.
func unauthorized(w http.ResponseWriter, token Token, realm string) {
	// Set WWW-Authenticate server response header for proper client auth
	if realm != "" {
		w.Header().Set("WWW-Authenticate", "Basic realm=\""+realm+":\"")
	}

	// Write unauthorized response
	w.WriteHeader(http.StatusUnauthorized)
	io.WriteString(w, "Unauthorized")
}
