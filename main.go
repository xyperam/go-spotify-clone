package main

import (
	"github.com/xyperam/go-spotify-clone/routes"
	"github.com/xyperam/go-spotify-clone/utils"
)

func main() {
	utils.ConnectDatabase()

	r := routes.SetupRoutes()
	r.Run(":8080")

}
