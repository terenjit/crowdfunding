package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Env struct {
	RootApp            string
	HTTPPort           uint16
	AccessTokenExpired time.Duration
	BasicAuthUsername  string
	BasicAuthPassword  string
	PostgreHost        string
	PostgreUser        string
	PostgrePassword    string
	PostgreDBName      string
	PostgrePort        uint16
	PostgreSSLMode     string
}

// GlobalEnv global environment
var GlobalEnv Env

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Println(err)
	}

	var ok bool

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	rootApp := strings.TrimSuffix(path, "/config")
	os.Setenv("APP_PATH", rootApp)
	GlobalEnv.RootApp = rootApp

	if port, err := strconv.Atoi(os.Getenv("HTTP_PORT")); err != nil {
		panic("missing HTTP_PORT environment")
	} else {
		GlobalEnv.HTTPPort = uint16(port)
	}

	GlobalEnv.BasicAuthUsername, ok = os.LookupEnv("BASIC_AUTH_USERNAME")
	if !ok {
		panic("missing BASIC_AUTH_USERNAME environment")
	}

	GlobalEnv.BasicAuthPassword, ok = os.LookupEnv("BASIC_AUTH_PASSWORD")
	if !ok {
		panic("missing BASIC_AUTH_PASSWORD environment")
	}

	GlobalEnv.PostgreHost, ok = os.LookupEnv("POSTGRE_HOST")
	if !ok {
		panic("missing POSTGRE_HOST environment")
	}

	GlobalEnv.PostgreUser, ok = os.LookupEnv("POSTGRE_USER")
	if !ok {
		panic("missing POSTGRE_USER environment")
	}

	GlobalEnv.PostgrePassword, ok = os.LookupEnv("POSTGRE_PASSWORD")
	if !ok {
		panic("missing POSTGRE_PASSWORD environment")
	}

	GlobalEnv.PostgreDBName, ok = os.LookupEnv("POSTGRE_DBNAME")
	if !ok {
		panic("missing POSTGRE_DBNAME environment")
	}

	if Portpostgre, err := strconv.Atoi(os.Getenv("POSTGRE_PORT")); err != nil {
		panic("missing POSTGRE_PORT environment")
	} else {
		GlobalEnv.PostgrePort = uint16(Portpostgre)
	}

	GlobalEnv.PostgreSSLMode, ok = os.LookupEnv("POSTGRE_SSLMODE")
	if !ok {
		panic("missing POSTGRE_SSLMODE environment")
	}

	GlobalEnv.AccessTokenExpired, err = time.ParseDuration("100m")
	if err != nil {
		panic("failed parsing AccessTokenExpired")
	}
}
