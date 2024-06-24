package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerAddress        string
	DatabaseDSN          string
	AccrualSystemAddress string
	JWTCookieName        string
	JWTSecret            string
	JWTTimeToLive        int
}

func New() Config {
	ParseFlags()

	return Config{
		ServerAddress:        getServerAddress(),
		DatabaseDSN:          getDatabaseDSN(),
		AccrualSystemAddress: getAccrualSystemAddress(),
		JWTCookieName:        getJWTCookieName(),
		JWTSecret:            getJWTSecret(),
		JWTTimeToLive:        getJWTTimeToLive(),
	}
}

func getServerAddress() string {
	var result = flagBind

	envServerAddress := os.Getenv("RUN_ADDRESS")
	if envServerAddress != "" {
		result = envServerAddress
	}

	return result
}

func getDatabaseDSN() string {
	var result = flagDatabaseDSN

	envDatabaseDSN := os.Getenv("DATABASE_URI")
	if envDatabaseDSN != "" {
		result = envDatabaseDSN
	}

	return result
}

func getAccrualSystemAddress() string {
	var result = accrualSystemAddress

	envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if envAccrualSystemAddress != "" {
		result = envAccrualSystemAddress
	}

	return result
}

func getJWTCookieName() string {
	var result = jwtCookieName

	envJwtCookieName := os.Getenv("JWT_COOKIE_NAME")
	if envJwtCookieName != "" {
		result = envJwtCookieName
	}

	return result
}

func getJWTSecret() string {
	var result = jwtSecret

	envJwtSecret := os.Getenv("JWT_SECRET")
	if envJwtSecret != "" {
		result = envJwtSecret
	}

	return result
}

func getJWTTimeToLive() int {
	var result = jwtTimeToLive

	envJwtTimeToLive := os.Getenv("JWT_TTL")
	if envJwtTimeToLive != "" {
		res, err := strconv.Atoi(envJwtTimeToLive)
		if err != nil {
			res = 0
		}
		result = res
	}

	return result
}
