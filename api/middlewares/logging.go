package middlewares

import (
	"log"
	"net/http"
	"time"
)

type Wrapper struct {
	http.ResponseWriter
	status int
}

func (w *Wrapper) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.status = code
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &Wrapper{w, http.StatusOK}
		next.ServeHTTP(wrapper, r)

		log.Println(wrapper.status, r.Method, r.URL.Path, time.Since(start))
	})
}
