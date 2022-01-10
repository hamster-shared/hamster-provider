package utils

import "math/rand"

func RandomPort() int {
	return 30000 + rand.Intn(10000)
}
