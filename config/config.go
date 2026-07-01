package config

import "os"

type EnvData struct{
	Port string
	DB string
	SessionSecret string
}

func Load() EnvData {
	 return EnvData{
		Port: os.Getenv("PORT"),
		DB: os.Getenv("DATABASE_URL"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
	 }
}