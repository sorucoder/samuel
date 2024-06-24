package server

import (
	"fmt"

	"github.com/sorucoder/samuel/internal/configuration"
)

func Run() {
	router := newRouter()

	address := fmt.Sprintf(`%s:%d`, configuration.Application.GetString("ip"), configuration.Application.GetInt("port"))
	fmt.Println(address)

	if configuration.Application.GetString("scheme") == "http" {
		router.Run(address)
	}
}
