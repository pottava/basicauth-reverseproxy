package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/justinas/alice"
)

func main() {

	// Proxy settings
	var proxyURL *url.URL
	if env := os.Getenv("PROXY_URL"); len(env) == 0 {
		log.Fatal("Missing required environment variables: PROXY_URL")
	} else {
		if parsed, err := url.Parse(env); err != nil {
			log.Fatalf("Could not parse proxy URL: %s", env)
		} else {
			log.Printf("[config] Proxy to %v", parsed)
			proxyURL = parsed
		}
	}
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.Director = func(r *http.Request) {
		r.Host = proxyURL.Host
		r.URL.Host = r.Host
		r.URL.Scheme = proxyURL.Scheme
	}
	http.Handle("/", alice.New(wrapper).Then(proxy))

	// Listen & Serve
	port := "80"
	if env := os.Getenv("APP_PORT"); len(env) > 0 {
		port = env
	}
	log.Printf("[service] listening on port %s", port)
	http.ListenAndServe(":"+port, nil)
}

func wrapper(h http.Handler) http.Handler {
	// Basic authentication
	user := os.Getenv("BASIC_AUTH_USER")
	pass := os.Getenv("BASIC_AUTH_PASS")
	if (len(user) > 0) && (len(pass) > 0) {
		log.Printf("[config] Basic authentication: %s", user)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (len(user) > 0) && (len(pass) > 0) && !auth(r, user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="REALM"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	})
}

func auth(r *http.Request, user, pass string) bool {
	if username, password, ok := r.BasicAuth(); ok {
		return username == user && password == pass
	}
	return false
}
