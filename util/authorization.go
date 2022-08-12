package util

import (
	"context"
	"errors"
	"os"

	"github.com/Nerzal/gocloak/v11"
)

const (
	AUTH_HEADER_PREFIX_LEN = 6
)

func AuthorizeUser(ctx context.Context, keycloak gocloak.GoCloak, authHeader string) (*gocloak.UserInfo, error) {
	if len(authHeader) < AUTH_HEADER_PREFIX_LEN {
		return nil, errors.New("invalid auth header")
	}

	authToken := authHeader[AUTH_HEADER_PREFIX_LEN:]
	tokenRetroRes, err := keycloak.RetrospectToken(ctx, authToken, os.Getenv("KEYCLOAK_CLIENT_ID"), os.Getenv("KEYCLOAK_CLIENT_SECRET"), os.Getenv("KEYCLOAK_REALM"))
	if err != nil {
		return nil, err
	}

	if tokenRetroRes.Active == nil || !(*tokenRetroRes.Active) {
		return nil, errors.New("token invalid")
	}

	return keycloak.GetUserInfo(ctx, authToken, os.Getenv("KEYCLOAK_REALM"))
}