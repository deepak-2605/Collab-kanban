package config

import "os"

type Config struct {
	MongoURI  string
	JWTSecret string
	Port      string
	DBName    string
}

func Load() Config {
	return Config{
		MongoURI:  os.Getenv("MONGO_URI"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      os.Getenv("PORT"),
		DBName:    "kanban",
	}
}
