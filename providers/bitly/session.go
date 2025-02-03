package bitly

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/haatos/goat"
)

// Session stores data during the auth process with Bitly.
type Session struct {
	AuthURL     string
	AccessToken string
}

// Ensure `bitly.Session` implements `goat.Session`.
var _ goat.Session = &Session{}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the Bitly provider.
func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goat.NoAuthUrlErrorMessage)
	}
	return s.AuthURL, nil
}

// Authorize the session with Bitly and return the access token to be stored for future use.
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
	return token.AccessToken, err
}

// Marshal the session into a string.
func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s Session) String() string {
	return s.Marshal()
}

// UnmarshalSession will unmarshal a JSON string into a session.
func (p *Provider) UnmarshalSession(data string) (goat.Session, error) {
	s := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(s)
	return s, err
}
