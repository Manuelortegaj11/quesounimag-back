package config

import (
  "github.com/joho/godotenv"
  "os"
  "log"
)


func InitEnv() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
}

func GetEnv(key string) string {
  return os.Getenv(key)
}
