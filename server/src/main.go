package main

import (
	"net/http"
	"playhouse-server/repository"
	"playhouse-server/router"
	"playhouse-server/util"
)

func main() {

	util.LoadEnv()
	f := repository.NewFactory()
	r := router.NewRootRouter(f)
	_ = http.ListenAndServe(":2345", r)
}
