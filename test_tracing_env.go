package main

import (
	"fmt"
	"github.com/weeb-vip/anime-api/config"
)

func main() {
	cfg := config.LoadConfigOrPanic()
	fmt.Println("Environment:", cfg.AppConfig.Env)
}