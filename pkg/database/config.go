package database

import (
	"net/url"
	"strconv"
)

type Config struct {
	Name              string `env:"DB_NAME, default=bike_station" json:",omitempty"`
	User              string `env:"DB_USER, default=bike_station" json:",omitempty"`
	Host              string `env:"DB_HOST, default=localhost" json:",omitempty"`
	Port              string `env:"DB_PORT, default=5432" json:",omitempty"`
	ConnectionTimeout int    `env:"DB_CONNECT_TIMEOUT" json:",omitempty"`
	Password          string `env:"DB_PASSWORD, default=secret" json:"-"`
	SSL               string `env:"DB_SSL_MODE"`
}

func (c *Config) DatabaseConfig() *Config {
	return c
}

func (c *Config) ConnectionURL() string {
	if c == nil {
		return ""
	}

	host := c.Host
	if v := c.Port; v != "" {
		host = host + ":" + v
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   c.Name,
	}

	if c.User != "" || c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}

	q := u.Query()
	if v := c.ConnectionTimeout; v > 0 {
		q.Add("connect_timeout", strconv.Itoa(v))
	}

	if v := c.SSL; len(v) > 0 {
		q.Add("sslmode", v)
	}

	u.RawQuery = q.Encode()

	return u.String()
}
