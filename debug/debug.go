package debug

import (
	"net/http"
	_ "net/http/pprof"
)

func init() {
	//noinspection GoUnhandledErrorResult
	go http.ListenAndServe(":6060", nil)
}
