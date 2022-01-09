// Author: InSun
// Date: 2021-01-09
// Description:
//	1. task 1 add input header to output header
//	2. task 2 set env VERSION to output header
//	3. task 3 write client_addr/statusCode/request_method into log
//	4. task 4 /healthz api with code 200
// 	5. task 5 elegant exit the http server
// Promotion（todo list）:
//	1. add config file for the server
//	2. reconstruct the project
// API:
// 	/healthz: health check
//	/test: test api
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test")
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	// task 4 healthz api with code 200
	w.WriteHeader(200)
}

func wrapHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// task 3
		log.Printf("--> %s %s from %s", r.Method, r.URL.Path, r.Host)
		lrw := &loggingResponseWriter{w, http.StatusOK}

		// task 1 add input header to output header
		for name, headers := range r.Header {
			for _, header := range headers {
				lrw.Header().Add(name, header)
			}
		}

		// task 2 set env VERSION to output header
		v := os.Getenv("VERSION")
		lrw.Header().Add("Version", v)

		// handler
		handler.ServeHTTP(lrw, r)

		// task 3
		statusCode := lrw.statusCode
		log.Printf("<-- %d %s", statusCode, http.StatusText(statusCode))
	})
}

func main() {
	var srv http.Server
	idleConnsClosed := make(chan struct{})

	// elegant exit http server
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	mux := http.NewServeMux()
	mux.Handle("/test", wrapHandler(http.HandlerFunc(testHandler)))
	mux.Handle("/healthz", wrapHandler(http.HandlerFunc(healthzHandler)))

	log.Printf("starting http server...")
	srv.Addr, srv.Handler = ":8090", mux
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
	// task 5 elegant exit server
	log.Printf("exit http server...")
}
