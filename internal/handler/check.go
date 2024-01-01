package handler

import "net/http"

func checkMethods(r *http.Request, methods []string) bool {
	for _, method := range methods {
		if r.Method == method {
			return true
		}
	}
	return false
}
