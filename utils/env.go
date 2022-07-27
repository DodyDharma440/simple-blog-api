package utils

import "os"

func GetEnv(key, fb string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fb
}
