package config

import "os"

type Config struct {
	DbName            string
	DbPort            string
	DbHost            string
	ServerPort        string
	PaginationPerPage string
}

func New() *Config {
	return &Config{
		DbName:            os.Getenv("DB_DATABASE"),
		DbPort:            os.Getenv("DB_PORT"),
		DbHost:            os.Getenv("DB_HOST"),
		ServerPort:        os.Getenv("SERVER_PORT"),
		PaginationPerPage: os.Getenv("PAGINATION_PER_PAGE"),
	}
}
