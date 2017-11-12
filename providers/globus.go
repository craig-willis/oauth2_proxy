package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type GlobusProvider struct {
	*ProviderData
}

//https://auth.globus.org/.well-known/openid-configuration
func NewGlobusProvider(p *ProviderData) *GlobusProvider {
	p.ProviderName = "Globus"
	if p.LoginURL == nil || p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{
			Scheme: "https",
			Host:   "auth.globus.org",
			Path:   "/v2/oauth2/authorize",
		}
	}
	if p.RedeemURL == nil || p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{
			Scheme: "https",
			Host:   "auth.globus.org",
			Path:   "/v2/oauth2/token",
		}
	}
	if p.ValidateURL == nil || p.ValidateURL.String() == "" {
		p.ValidateURL = &url.URL{
			Scheme: "https",
			Host:   "auth.globus.org",
			Path:   "/v2/oauth2/token/introspect",
		}
	}
	if p.ProfileURL == nil || p.ProfileURL.String() == "" {
		p.ProfileURL = &url.URL{
			Scheme: "https",
			Host:   "auth.globus.org",
			Path:   "/v2/oauth2/userinfo",
		}
	}
	if p.Scope == "" {
		p.Scope = "openid email profile"
	}
	return &GlobusProvider{ProviderData: p}
}

func (p *GlobusProvider) GetEmailAddress(s *SessionState) (string, error) {
	var userinfo struct {
		Email             string `json:"email"`
		Name              string `json:"name"`
		PreferredUsername string `json:"preferred_username"`
		Sub               string `json:"sub"`
	}

	req, _ := http.NewRequest("GET", p.ProfileURL.String(), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("got %d from %q %s",
			resp.StatusCode, p.ProfileURL.String(), body)
	} else {
		log.Printf("got %d from %q %s", resp.StatusCode, p.ProfileURL.String(), body)
	}

	if err := json.Unmarshal(body, &userinfo); err != nil {
		return "", fmt.Errorf("%s unmarshaling %s", err, body)
	}

	return userinfo.Email, nil
}
