package config

import "os"

type EnvData struct{
	port string
	DB string
	sessionSecret string
}

func Load() EnvData {
	 return EnvData{
		port: os.Getenv("PORT"),
		DB: os.Getenv("DATABASE_URL"),
		sessionSecret: os.Getenv("SESSON_SECRET"),
	 }
}