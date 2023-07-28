package main

import (
	"net/http"
	"playhouse-server/env"
	"playhouse-server/router"
)

func main() {
	env.Load()
	r := router.NewRootRouter()
	_ = http.ListenAndServe(":2345", r)
}
