package config

import (
	"flag"
	"fmt"
)

const defaultPort = 8080

var flagBind = fmt.Sprintf(":%d", defaultPort)
var flagDatabaseDSN = "postgres://postgres:mysecretpassword@127.0.0.1:5432/postgres"
var accrualSystemAddress = "???"

func parseFlags() {
	flag.StringVar(&flagBind, "a", flagBind, "Host and port to bind")
	flag.StringVar(&flagDatabaseDSN, "d", flagDatabaseDSN, "Database DSN")
	flag.StringVar(&accrualSystemAddress, "r", accrualSystemAddress, "Accrual system address")

	flag.Parse()
}
