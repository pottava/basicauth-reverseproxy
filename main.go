package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"github.com/justinas/alice"
)

type config struct {
	proxyURL      *url.URL // PROXY_URL
	basicAuthUser string   // BASIC_AUTH_USER
	basicAuthPass string   // BASIC_AUTH_PASS
	port          string   // APP_PORT
	accessLog     bool     // ACCESS_LOG
	sslCert       string   // SSL_CERT_PATH
	sslKey        string   // SSL_KEY_PATH
}

var (
	version string
	date    string
	c       *config
)

func main() {
	c = configFromEnvironmentVariables()
	proxy := httputil.NewSingleHostReverseProxy(c.proxyURL)
	proxy.Director = func(r *http.Request) {
		r.Host = c.proxyURL.Host
		r.URL.Host = r.Host
		r.URL.Scheme = c.proxyURL.Scheme
	}
	http.Handle("/", alice.New(wrapper).Then(proxy))

	// Listen & Serve
	log.Printf("[service] listening on port %s", c.port)
	if (len(c.sslCert) > 0) && (len(c.sslKey) > 0) {
		log.Fatal(http.ListenAndServeTLS(":"+c.port, c.sslCert, c.sslKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+c.port, nil))
	}
}

func configFromEnvironmentVariables() *config {
	candidate := os.Getenv("PROXY_URL")
	if len(candidate) == 0 {
		log.Fatal("Missing required environment variable: PROXY_URL")
	}
	proxyURL, err := url.Parse(candidate)
	if err != nil {
		log.Fatalf("Could not parse proxy URL: %s", candidate)
	}
	port := os.Getenv("APP_PORT")
	if len(port) == 0 {
		port = "80"
	}
	accessLog := false
	if candidate, found := os.LookupEnv("ACCESS_LOG"); found {
		if b, err := strconv.ParseBool(candidate); err == nil {
			accessLog = b
		}
	}
	conf := &config{
		proxyURL:      proxyURL,
		basicAuthUser: os.Getenv("BASIC_AUTH_USER"),
		basicAuthPass: os.Getenv("BASIC_AUTH_PASS"),
		port:          port,
		accessLog:     accessLog,
		sslCert:       os.Getenv("SSL_CERT_PATH"),
		sslKey:        os.Getenv("SSL_KEY_PATH"),
	}
	// Proxy
	log.Printf("[config] Proxy to %v", proxyURL)

	// TLS pem files
	if (len(conf.sslCert) > 0) && (len(conf.sslKey) > 0) {
		log.Print("[config] TLS enabled.")
	}
	// Basic authentication
	if (len(conf.basicAuthUser) > 0) && (len(conf.basicAuthPass) > 0) {
		log.Printf("[config] Basic authentication: %s", conf.basicAuthUser)
	}
	return conf
}

func wrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (len(c.basicAuthUser) > 0) && (len(c.basicAuthPass) > 0) && !auth(r, c.basicAuthUser, c.basicAuthPass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="REALM"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
		if c.accessLog {
			log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		}
	})
}

func auth(r *http.Request, user, pass string) bool {
	if username, password, ok := r.BasicAuth(); ok {
		return username == user && password == pass
	}
	return false
}
