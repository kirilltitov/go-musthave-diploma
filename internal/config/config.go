package config

import (
	"os"
)

type Config struct {
	ServerAddress        string
	DatabaseDSN          string
	AccrualSystemAddress string
}

func New() Config {
	parseFlags()

	return Config{
		ServerAddress:        getServerAddress(),
		DatabaseDSN:          getDatabaseDSN(),
		AccrualSystemAddress: getAccrualSystemAddress(),
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
