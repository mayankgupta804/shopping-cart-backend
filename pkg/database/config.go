package database

import (
	"net/url"
	"strconv"
)

type Config struct {
	Name              string
	User              string
	Host              string
	Port              string
	Password          string
	SSL               string
	ConnectionTimeout int
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
