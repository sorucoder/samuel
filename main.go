package main

import (
	"github.com/sorucoder/samuel/internal/configuration"
	"github.com/sorucoder/samuel/internal/database"
	"github.com/sorucoder/samuel/internal/email"
	"github.com/sorucoder/samuel/internal/ldap"
	"github.com/sorucoder/samuel/internal/server"
)

func main() {
	configuration.Initialize()
	database.Initialize()
	ldap.Initialize()
	email.Initialize()

	server.Run()
}
