package middleware

import (
	"log"
	"net/http"
	"time"
)

// Middleware is a custom middleware
type Middleware struct {
}

// LoggingHandler logs the http request time
func (m Middleware) LoggingHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//FIXME: broken
	t1 := time.Now()
	next.ServeHTTP(rw, r)
	t2 := time.Now()
	log.Printf("[%s] %q %v", r.Method, r.URL.String(), t2.Sub(t1))
	next(rw, r)
}

// RecoverHandler recover from panic
func (m Middleware) RecoverHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//FIXME: broken
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recover from panic: %+v", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()
	next.ServeHTTP(rw, r)
}
