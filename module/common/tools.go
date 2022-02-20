package common

import (
	"os"
)

var (
	Done = make(chan struct{})

	EnvType string
)

func GetEnvType() string {
	envType, ok := os.LookupEnv("EnvType")
	if !ok || envType == "" {
		EnvType = "Testing"
	} else {
		EnvType = envType
	}

	return EnvType
}
