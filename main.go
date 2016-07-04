package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type proxyPattern struct {
	matches  []string
	proxyURL *url.URL
}

type config struct {
	ProxyPatterns []proxyPattern
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

	// Proxy!!
	dummy, _ := url.Parse("https://www.docker.com/")
	proxy := httputil.NewSingleHostReverseProxy(dummy)
	proxy.Director = func(r *http.Request) {
		r.Header.Set("Host", r.Host)

		found := false
		for _, patterns := range c.ProxyPatterns {
			if match(patterns.matches[0], r.URL.Scheme, false) &&
				match(patterns.matches[1], r.Host, false) &&
				match(patterns.matches[3], r.URL.Path, true) {
				r.Host = patterns.proxyURL.Host
				r.URL.Host = r.Host
				r.URL.Scheme = patterns.proxyURL.Scheme
				found = true
				break
			}
		}
		if !found {
			r.Host = c.proxyURL.Host
			r.URL.Host = r.Host
			r.URL.Scheme = c.proxyURL.Scheme
		}
		if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			if prior, ok := r.Header["X-Forwarded-For"]; ok {
				ip = strings.Join(prior, ", ") + ", " + ip
			}
			r.Header.Set("X-Forwarded-For", ip)
		}
	}
	http.Handle("/", wrapper(proxy))

	http.HandleFunc("/--version", func(w http.ResponseWriter, r *http.Request) {
		if len(version) > 0 && len(date) > 0 {
			fmt.Fprintf(w, "version: %s (built at %s)", version, date)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	// Listen & Serve
	log.Printf("[service] listening on port %s", c.port)
	if (len(c.sslCert) > 0) && (len(c.sslKey) > 0) {
		log.Fatal(http.ListenAndServeTLS(":"+c.port, c.sslCert, c.sslKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+c.port, nil))
	}
}

func match(pattern, target string, isPath bool) bool {
	if pattern == "" {
		return true
	}
	pattern = strings.Replace(pattern, "*", ".*", -1)
	if isPath && !strings.HasSuffix(pattern, "/") {
		pattern += "$"
	}
	match, _ := regexp.MatchString(pattern, target)
	return match
}

func configFromEnvironmentVariables() *config {
	candidateProxyURL := strings.Trim(os.Getenv("PROXY_URL"), "\"")
	candidateProxyPatterns := strings.Trim(os.Getenv("PROXY_PATTERNS"), "\"")
	if len(candidateProxyURL) == 0 && len(candidateProxyPatterns) == 0 {
		log.Fatal("Missing required environment variable: PROXY_URL or PROXY_PATTERNS")
	}
	var proxyURL *url.URL
	var err error
	if len(candidateProxyURL) > 0 {
		proxyURL, err = url.Parse(candidateProxyURL)
		if err != nil {
			log.Fatalf("Could not parse proxy URL: %s", candidateProxyURL)
		}
	}
	ProxyPatterns := []proxyPattern{}
	if len(candidateProxyPatterns) > 0 {
		regex := regexp.MustCompile(`(https?://)?([^:^/]*)(:\\d*)?(.*)?`)
		for _, pattern := range strings.Split(candidateProxyPatterns, ",") {
			splited := strings.Split(pattern, "=")
			if url, err := url.Parse(splited[1]); err == nil {
				pattern := proxyPattern{
					matches:  regex.FindStringSubmatch(splited[0])[1:],
					proxyURL: url,
				}
				ProxyPatterns = append(ProxyPatterns, pattern)
			}
		}
	}
	port := os.Getenv("APP_PORT")
	if len(port) == 0 {
		port = "80"
	}
	accessLog := false
	if b, err := strconv.ParseBool(os.Getenv("ACCESS_LOG")); err == nil {
		accessLog = b
	}
	conf := &config{
		ProxyPatterns: ProxyPatterns,
		proxyURL:      proxyURL,
		basicAuthUser: os.Getenv("BASIC_AUTH_USER"),
		basicAuthPass: os.Getenv("BASIC_AUTH_PASS"),
		port:          port,
		accessLog:     accessLog,
		sslCert:       os.Getenv("SSL_CERT_PATH"),
		sslKey:        os.Getenv("SSL_KEY_PATH"),
	}
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
