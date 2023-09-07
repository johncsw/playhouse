package main

import (
	"net/http"
	"playhouse-server/env"
	"playhouse-server/flow"
	"playhouse-server/repo"
	"playhouse-server/router"
)

func main() {
	env.Load()
	repo.SetUpSchema()
	repo.Init()
	flow.DeleteAllDataEveryHour()
	r := router.NewRootRouter()
	_ = http.ListenAndServe(":2345", r)
}
