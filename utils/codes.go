package utils

import (
    "math/rand"
    "time"
    "fmt"
)

func GenerateConfirmationCode() string {
    source := rand.NewSource(time.Now().UnixNano())
    r := rand.New(source)
    code := r.Intn(1000000)
    return fmt.Sprintf("%06d", code)
}
