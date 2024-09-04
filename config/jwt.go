package config

type Jwt struct {
	TokenExpire int64  `mapstructure:"token-expire" json:"token-expire" yaml:"token-expire"`
	SigningKey  string `mapstructure:"signing-key" json:"signing-key" yaml:"signing-key"`
}
