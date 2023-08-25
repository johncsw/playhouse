package main

import (
	"net/http"
	"playhouse-server/env"
	"playhouse-server/repo"
	"playhouse-server/router"
)

func main() {
	env.Load()
	repo.Init()
	r := router.NewRootRouter()
	_ = http.ListenAndServe(":2345", r)
}
