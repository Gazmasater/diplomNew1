package main

import (
	"net/http"

	"diplom.com/internal/config"
	"diplom.com/internal/router"
)

// main
func main() {
	cfg := config.Load()

	r := router.NewRouter(cfg)

	http.ListenAndServe(cfg.Address, r)
}
