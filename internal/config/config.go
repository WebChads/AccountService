package config

type ServerConfig struct {
	port  string `yaml:"port"`
	dbUrl string `yaml:"database_url"`
}

func NewServerConfig() *ServerConfig {
	
}
