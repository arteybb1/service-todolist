package main

import (
	"github.com/arteybb/service-todolist/internal/config"
	"github.com/arteybb/service-todolist/internal/router"
)

func main() {
	config.LoadConfig()
	config.MongoConfig()
	r := router.Route()
	r.Run(":" + config.AppConfig.Port)
}
