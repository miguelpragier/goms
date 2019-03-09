package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func mwJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("authenticating access for %v", r.URL.Path)

		if jwt.CheckHeaderAndRenew(w, r) {
			next(w, r)
			return
		}

		log.Println("not authorized")

		http.Error(w, "Acesso não autorizado. Reexecute a autenticação", http.StatusUnauthorized)
	}
}

func mwLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s %s", r.Method, r.URL.Path, r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func webserverStart() {
	router := mux.NewRouter()

	// Heartbeat/ping test
	router.HandleFunc("/api/ping/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("PONG!!"))
	}).Methods(http.MethodGet)

	// Running time check
	router.HandleFunc("/api/uptime/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(uptime()))
	}).Methods(http.MethodGet)

	// Retrieve gitRevision
	router.HandleFunc("/api/gitrevision/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(gitRevisionHash))
	}).Methods(http.MethodGet)

	// Retrieve compilation
	// Reference: https://blog.alexellis.io/inject-build-time-vars-golang/
	router.HandleFunc("/api/compilationtimestamp/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(compilationTimestamp))
	}).Methods(http.MethodGet)

	http.Handle(`/`, router)

	log.Println("Server listening on port ", webserverListeningPort)

	s := fmt.Sprintf(":%d", webserverListeningPort)

	//router.Use(mwOpenCORS, mwLog)

	allowedOrigin := os.Getenv("allowedOrigin")

	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	originsOk := handlers.AllowedOrigins([]string{allowedOrigin})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Origin"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	hnd := handlers.CORS(headersOk, originsOk, methodsOk, handlers.AllowCredentials())(router)

	log.Fatal(http.ListenAndServe(s, hnd))
}
