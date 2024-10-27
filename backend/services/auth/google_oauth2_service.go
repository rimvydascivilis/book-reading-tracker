package auth

import (
	"context"
	"fmt"

	"google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
)

type GoogleOAuth2Service struct {
	oauth2Svc *oauth2.Service
}

func NewGoogleOAuth2Service() (*GoogleOAuth2Service, error) {
	ctx := context.Background()
	oauthSvc, err := oauth2.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}

	return &GoogleOAuth2Service{
		oauth2Svc: oauthSvc,
	}, nil
}

func (g *GoogleOAuth2Service) ValidateToken(token string) (string, error) {
	tokenInfoCall := g.oauth2Svc.Tokeninfo()
	tokenInfoCall.IdToken(token)

	tokenInfo, err := tokenInfoCall.Do()
	if err != nil || tokenInfo.Email == "" {
		return "", fmt.Errorf("failed to validate OAuth token: %w", err)
	}

	return tokenInfo.Email, nil
}
