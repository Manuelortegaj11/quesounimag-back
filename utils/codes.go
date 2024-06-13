package utils

import (
    "math/rand"
    "time"
    "fmt"
)

func GenerateConfirmationCode() string {
    rand.Seed(time.Now().UnixNano())
    code := rand.Intn(1000000)
    return fmt.Sprintf("%06d", code)
}
