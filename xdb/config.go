package xdb

// Config the xdb config
type Config struct {
	Dsn         string `yaml:"dsn"`
	MaxIdle     int    `yaml:"maxIdle"`
	MaxOpen     int    `yaml:"maxOpen"`
	MaxLifetime int    `yaml:"maxLifetime"`
	LogMode     bool   `yaml:"logMode"`
}
