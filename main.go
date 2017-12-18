package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
)

type Proxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewProxy(target string) (*Proxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		target: u,
		proxy:  httputil.NewSingleHostReverseProxy(u),
	}, nil
}

func (p *Proxy) handle(w http.ResponseWriter, r *http.Request) {
	local, _ := url.Parse(p.target.String())
	local.Path = path.Join(local.Path, r.URL.Path)

	log.Printf("%s -> %s", r.Method, local.String())

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")

	r.Host = r.URL.Host
	p.proxy.ServeHTTP(w, r)
}

func envdef(key, defval string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}

	return defval
}

func main() {
	var (
		defaultPort        = envdef("CORS_PORT", ":8889")
		defaultPortUsage   = "The default port for the reverse CORS proxy"
		defaultTarget      = envdef("CORS_URL", "http://localhost/")
		defaultTargetUsage = "The default URL to use when doing reverse proxy"
	)

	port := flag.String("port", defaultPort, defaultPortUsage)
	url := flag.String("url", defaultTarget, defaultTargetUsage)

	flag.Parse()

	log.Printf("Server listening on: %q for URL: %s", *port, *url)

	proxy, err := NewProxy(*url)
	if err != nil {
		log.Fatalf("Unable to parse proxied URL: %s", err.Error())
	}

	http.HandleFunc("/", proxy.handle)
	if err := http.ListenAndServe(*port, nil); err != nil {
		log.Fatalf("Error while launching the HTTP server: %s", err.Error())
	}
}
