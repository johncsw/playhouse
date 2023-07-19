package main

import (
	"net/http"
	"playhouse-server/router"
	"playhouse-server/util"
)

func main() {

	env := util.NewEnv()
	env.Load()
	r := router.NewRootRouter()
	_ = http.ListenAndServe(":2345", r)
}
