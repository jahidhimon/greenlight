package greenlog

import (
	"log"
	"net/http"
)

func Logreq(r *http.Request, path string) {
	header := r.Header
	agent := header.Get("User-Agent")
	method := r.Method
	log.Print(agent, " (", method, " ", path, ")")
}
