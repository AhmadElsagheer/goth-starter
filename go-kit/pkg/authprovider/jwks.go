package authprovider

import (
	"context"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type JWKSProvider struct {
	jwks keyfunc.Keyfunc
}

func NewJWKSProvider(jwksURL string) (Provider, error) {
	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS from resource at the given URL: %w", err)
	}
	return &JWKSProvider{
		jwks: jwks,
	}, nil
}

func (p *JWKSProvider) VerifyIDToken(ctx context.Context, tokenString string) (IDToken, error) {
	token, err := jwt.Parse(tokenString, p.jwks.Keyfunc)
	if err != nil {
		return IDToken{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return IDToken{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return IDToken{}, fmt.Errorf("invalid token claims")
	}

	// Verify expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return IDToken{}, fmt.Errorf("token expired")
		}
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return IDToken{}, fmt.Errorf("sub claim not found in token")
	}

	return IDToken{
		ID: sub,
	}, nil
}
