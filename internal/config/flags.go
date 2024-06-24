package config

import (
	"flag"
	"fmt"
)

const defaultPort = 8080

var flagBind = fmt.Sprintf(":%d", defaultPort)
var flagDatabaseDSN = "postgres://postgres:mysecretpassword@127.0.0.1:5432/postgres"
var accrualSystemAddress = "???"
var jwtCookieName = "access_token"
var jwtSecret = "hesoyam"
var jwtTimeToLive int = 86400

func ParseFlags() {
	flag.StringVar(&flagBind, "a", flagBind, "Host and port to bind")
	flag.StringVar(&flagDatabaseDSN, "d", flagDatabaseDSN, "Database DSN")
	flag.StringVar(&accrualSystemAddress, "r", accrualSystemAddress, "Accrual system address")
	flag.StringVar(&jwtCookieName, "cookie-name", jwtCookieName, "JWT Cookie name")
	flag.StringVar(&jwtSecret, "jwt-secret", jwtSecret, "JWT Secret")
	flag.IntVar(&jwtTimeToLive, "jwt-ttl", jwtTimeToLive, "JWT Time To Live")

	flag.Parse()
}
