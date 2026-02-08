package internal

import "time"

type PostgresConfig struct {
	Database string         `koanf:"database"`
	Schema   string         `koanf:"schema"`
	Config   DatabaseConfig `koanf:"config"`
}

type DatabaseSSL struct {
	Enabled bool `koanf:"enabled"`
}

type DatabaseConfig struct {
	Host    string        `koanf:"host"`
	Port    uint16        `koanf:"port"`
	Auth    DatabaseAuth  `koanf:"auth"`
	SSL     DatabaseSSL   `koanf:"ssl"`
	Timeout time.Duration `koanf:"timeout"`
	Pool    DatabasePool  `koanf:"pool"`
}

type DatabaseAuth struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

type DatabasePool struct {
	MinConns        int32         `koanf:"min"`
	MaxConns        int32         `koanf:"max"`
	MaxConnLifetime time.Duration `koanf:"lifetime"`
	MaxConnIdleTime time.Duration `koanf:"idletime"`
}
