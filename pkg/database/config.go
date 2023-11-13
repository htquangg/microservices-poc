package database

import (
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	Port            uint16
	Host            string
	User            string
	Password        string
	Schema          string
	Charset         string
	AutoMigration   bool
	LogSQL          bool
	SslMode         bool
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

func (c *Config) Address() string {
	param := "?"
	if strings.Contains(c.Schema, param) {
		param = "&"
	}

	conn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v%scharset=%s&parseTime=true&tls=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Schema,
		param,
		c.Charset,
		strconv.FormatBool(c.SslMode),
	)

	return conn
}
