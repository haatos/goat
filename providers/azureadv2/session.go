package azureadv2

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/haatos/goat"
)

// Session is the implementation of `goat.Session`
type Session struct {
	AuthURL      string    `json:"au"`
	AccessToken  string    `json:"at"`
	RefreshToken string    `json:"rt"`
	ExpiresAt    time.Time `json:"exp"`
}

// GetAuthURL will return the URL set by calling the `BeginAuth` func
func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goat.NoAuthUrlErrorMessage)
	}

	return s.AuthURL, nil
}

// Authorize the session with AzureAD and return the access token to be stored for future use.
func (s *Session) Authorize(provider goat.Provider, params goat.Params) (string, error) {
	p := provider.(*Provider)
	token, err := p.config.Exchange(goat.ContextForClient(p.Client()), params.Get("code"))
	if err != nil {
		return "", err
	}

	if !token.Valid() {
		return "", errors.New("invalid token received from provider")
	}

	s.AccessToken = token.AccessToken
	s.RefreshToken = token.RefreshToken
	s.ExpiresAt = token.Expiry

	return token.AccessToken, err
}

// Marshal the session into a string
func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s Session) String() string {
	return s.Marshal()
}

// UnmarshalSession wil unmarshal a JSON string into a session.
func (p *Provider) UnmarshalSession(data string) (goat.Session, error) {
	session := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(session)
	return session, err
}
