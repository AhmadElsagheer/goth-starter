package authprovider

import "context"

//go:generate mockgen -source ./authprovider.go -destination ./mockauthprovider/mockauthprovider.go -package mockauthprovider
type Provider interface {
	VerifyIDToken(ctx context.Context, token string) (IDToken, error)
}

type IDToken struct {
	ID string
}
