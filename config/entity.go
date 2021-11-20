package config

type Config struct {
	Env       string   `mapstructure:"env"`
	Port      string   `mapstructure:"port"`
	JwtSecret string   `mapstructure:"jwtSecret"`
	JwtIssuer string   `mapstructure:"jwtIssuer"`
	Database  Database `mapstructure:"database"`
}

type Database struct {
	DBName   string `mapstructure:"dbName"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}
