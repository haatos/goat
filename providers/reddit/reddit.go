package reddit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/haatos/goat"
	"golang.org/x/oauth2"
)

const (
	authURL = "https://www.reddit.com/api/v1/authorize"
)

type Provider struct {
	providerName string
	duration     string
	config       oauth2.Config
	client       http.Client
	// TODO: userURL should be a constant
	userURL string
}

func New(
	clientID string,
	clientSecret string,
	redirectURI string,
	duration string,
	tokenEndpoint string,
	userURL string,
	scopes ...string,
) Provider {
	return Provider{
		providerName: "reddit",
		duration:     duration,
		config: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   authURL,
				TokenURL:  tokenEndpoint,
				AuthStyle: 0,
			},
			RedirectURL: redirectURI,
			Scopes:      scopes,
		},
		client:  http.Client{},
		userURL: userURL,
	}
}

func (p *Provider) Name() string {
	return p.providerName
}

func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) UnmarshalSession(s string) (goat.Session, error) {
	session := &Session{}
	err := json.Unmarshal([]byte(s), session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (p *Provider) Debug(b bool) {}

func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return nil, nil
}

func (p *Provider) RefreshTokenAvailable() bool {
	return true
}

func (p *Provider) BeginAuth(state string) (goat.Session, error) {
	authCodeOption := oauth2.SetAuthURLParam("duration", p.duration)
	return &Session{AuthURL: p.config.AuthCodeURL(state, authCodeOption)}, nil
}

type redditResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (p *Provider) FetchUser(s goat.Session) (goat.User, error) {
	session := s.(*Session)
	request, err := http.NewRequest("GET", p.userURL, nil)
	if err != nil {
		return goat.User{}, err
	}

	bearer := "Bearer " + session.AccessToken
	request.Header.Add("Authorization", bearer)

	res, err := p.client.Do(request)
	if err != nil {
		return goat.User{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusForbidden {
			return goat.User{}, fmt.Errorf(
				"%s responded with a %s because you did not provide the identity scope which is required to fetch user profile",
				p.providerName,
				res.Status,
			)
		}
		return goat.User{}, fmt.Errorf(
			"%s responded with a %d trying to fetch user profile",
			p.providerName,
			res.StatusCode,
		)
	}

	bits, err := io.ReadAll(res.Body)
	if err != nil {
		return goat.User{}, err
	}

	var r redditResponse

	err = json.Unmarshal(bits, &r)
	if err != nil {
		return goat.User{}, err
	}

	goatUser := goat.User{
		RawData:      nil,
		Provider:     p.Name(),
		Name:         r.Name,
		UserID:       r.Id,
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    time.Time{},
	}

	err = json.Unmarshal(bits, &goatUser.RawData)
	if err != nil {
		return goat.User{}, err
	}

	return goatUser, nil
}
