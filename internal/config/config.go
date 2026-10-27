package config

import (
	"net"
	"net/url"
)

type Config struct {
	App Application
	Db  Database
}

type Application struct {
	Name string
	Port string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (c *Config) DatabaseURL() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.Db.User, c.Db.Password),
		Host:   net.JoinHostPort(c.Db.Host, c.Db.Port),
		Path:   c.Db.Name,
	}

	q := u.Query()
	q.Set("sslmode", c.Db.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}
