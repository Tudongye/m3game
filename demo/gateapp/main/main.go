package main

import (
	"context"
	"m3game/demo/gateapp"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go http.ListenAndServe(":9999", nil)
	gateapp.Run(context.Background())
}
