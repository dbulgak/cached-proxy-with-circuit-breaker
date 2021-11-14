package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if len(value) == 0 {
		return defaultValue
	}

	return value
}

func GetIntEnv(key string, defaultValue int) int {
	valuestr := os.Getenv(key)

	if len(valuestr) == 0 {
		log.Errorf("no %s value, setting default %d", key, defaultValue)
		return defaultValue
	}

	valueint, err := strconv.Atoi(valuestr)

	if err != nil {
		log.Errorf("unexpected conversion of %s, %s, using default valuestr", valuestr, err)
		return defaultValue
	}

	return valueint
}
