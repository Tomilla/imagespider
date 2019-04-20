package debug

import (
	"net/http"
	_ "net/http/pprof"
)

func init() {
	go http.ListenAndServe(":6060", nil)
}
