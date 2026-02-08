package authprovider

type ProviderType string

const (
	JWTProviderType ProviderType = "jwt"
)

type JWTProviderConfig struct {
	PublicKey string `koanf:"publickey"`
}
