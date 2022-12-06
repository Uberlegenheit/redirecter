package conf

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DefaultRedirectURL      string
	RedirectHost            string
	CORSAllowedOrigins      []string
	ShortcutterListenOnPort uint64
	Postgres                Postgres
	WaitTimeout             uint64
}

type Postgres struct {
	Host     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func NewFromENV() (Config, error) {
	shortcutPort, err := strconv.ParseUint(os.Getenv("LISTEN_PORT"), 10, 64)
	waitTimeout, err := strconv.ParseUint(os.Getenv("WAIT_TIMEOUT"), 10, 64)
	if err != nil {
		return Config{}, err
	}
	return Config{
		DefaultRedirectURL:      os.Getenv("DEFAULT_REDIRECT_URL"),
		RedirectHost:            os.Getenv("REDIRECT_HOST"),
		CORSAllowedOrigins:      strings.Split(os.Getenv("CORS_ALLOWED"), ","),
		ShortcutterListenOnPort: shortcutPort,
		WaitTimeout:             waitTimeout,
		Postgres: Postgres{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSL"),
		},
	}, nil
}
