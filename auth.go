package main

import "net/http"

type AuthMiddleware struct {
	handler http.Handler
}

func (am *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("key") != CFG.Key {
		w.WriteHeader(401)
		return
	}
	am.handler.ServeHTTP(w, r)
}
