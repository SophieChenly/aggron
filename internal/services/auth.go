package services

import (
	"context"
	"log"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type AuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	IssuerURL    string
}

type AuthService interface {
	Authorize(state string, authCodeOpt oauth2.AuthCodeOption) string
}

type Auth struct {
	Config oauth2.Config
}

func NewAuth(config AuthConfig) (*Auth, error) {
	provider, err := oidc.NewProvider(context.Background(), config.IssuerURL)
	if err != nil {
		log.Fatalf("Failed to create OIDC provider: %v", err)
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       config.Scopes,
	}

	return &Auth{
		Config: oauth2Config,
	}, nil
}

func (a *Auth) Authorize(state string, authCodeOpt oauth2.AuthCodeOption) string {
	url := a.Config.AuthCodeURL(state, authCodeOpt)

	return url
}
