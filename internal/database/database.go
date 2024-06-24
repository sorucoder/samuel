package database

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sorucoder/samuel/internal/configuration"
)

var (
	driver string
	dsn    string
	raw    *sqlx.DB
)

var (
	ErrUnsupportedDriver error = errors.New("database unsupported driver")
)

func Initialize() {
	driver = configuration.Database.GetString("driver")

	dsn = fmt.Sprintf(
		`%s:%s@tcp(%s:%d)/samuel?parseTime=true&loc=Local`,
		configuration.Database.GetString("user"),
		configuration.Database.GetString("password"),
		configuration.Database.GetString("host"),
		configuration.Database.GetInt("port"),
	)

	var errOpen error
	raw, errOpen = sqlx.Open(driver, dsn)
	if errOpen != nil {
		panic(errOpen)
	}
}

func Ping(context context.Context) error {
	return raw.PingContext(context)
}
