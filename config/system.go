package config

type System struct {
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	UseRedis bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`
}
