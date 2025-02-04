package microsoftonline

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/haatos/goat"
)

// Session is the implementation of `goat.Session` for accessing microsoftonline.
// Refresh token not available for microsoft online: session size hit the limit of max cookie size
type Session struct {
	AuthURL     string
	AccessToken string
	ExpiresAt   time.Time
}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the Facebook provider.
func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goat.NoAuthUrlErrorMessage)
	}

	return s.AuthURL, nil
}

// Authorize the session with Facebook and return the access token to be stored for future use.
func (s *Session) Authorize(provider goat.Provider, params goat.Params) (string, error) {
	p := provider.(*Provider)
	token, err := p.config.Exchange(goat.ContextForClient(p.Client()), params.Get("code"))
	if err != nil {
		return "", err
	}

	if !token.Valid() {
		return "", errors.New("Invalid token received from provider")
	}

	s.AccessToken = token.AccessToken
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
