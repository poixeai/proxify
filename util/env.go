package util

import "os"

func GetEnvPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7777"
	}

	return port
}
