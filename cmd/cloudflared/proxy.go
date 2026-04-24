package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// ProxyCommand defines the CLI command for running a local proxy server
// that forwards traffic through the Cloudflare tunnel.
var ProxyCommand = &cli.Command{
	Name:      "proxy",
	Usage:     "Run a local proxy server to forward traffic through the tunnel",
	ArgsUsage: "[origin-url]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"a"},
			Usage:   "Address to bind the proxy server on",
			Value:   "127.0.0.1",
		},
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "Port to bind the proxy server on",
			Value:   8080,
		},
		&cli.DurationFlag{
			Name:  "timeout",
			Usage: "Timeout for upstream requests",
			Value: 30 * time.Second,
		},
		&cli.BoolFlag{
			Name:  "tls-verify",
			Usage: "Verify TLS certificates of the origin server",
			Value: true,
		},
	},
	Action: runProxy,
}

// ProxyConfig holds configuration for the local proxy server.
type ProxyConfig struct {
	ListenAddr string
	OriginURL  *url.URL
	Timeout    time.Duration
}

// runProxy starts the local proxy server with the provided CLI context.
func runProxy(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("origin URL is required")
	}

	originRaw := c.Args().First()
	originURL, err := url.Parse(originRaw)
	if err != nil {
		return fmt.Errorf("invalid origin URL %q: %w", originRaw, err)
	}

	if originURL.Scheme == "" {
		originURL.Scheme = "http"
	}

	addr := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))
	cfg := &ProxyConfig{
		ListenAddr: addr,
		OriginURL:  originURL,
		Timeout:    c.Duration("timeout"),
	}

	return startProxyServer(cfg)
}

// startProxyServer initializes and starts the HTTP proxy server.
func startProxyServer(cfg *ProxyConfig) error {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   cfg.Timeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: cfg.Timeout,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
	}

	handler := &proxyHandler{
		origin:    cfg.OriginURL,
		client:    &http.Client{Transport: transport, Timeout: cfg.Timeout},
	}

	server := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      handler,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}

	log.Info().Str("addr", cfg.ListenAddr).Str("origin", cfg.OriginURL.String()).Msg("Starting proxy server")
	return server.ListenAndServe()
}

// proxyHandler implements http.Handler and forwards requests to the origin.
type proxyHandler struct {
	origin *url.URL
	client *http.Client
}

// ServeHTTP proxies the incoming request to the configured origin URL.
func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := *h.origin
	target.Path = r.URL.Path
	target.RawQuery = r.URL.RawQuery

	req, err := http.NewRequestWithContext(r.Context(), r.Method, target.String(), r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create upstream request")
		http.Error(w, "proxy error", http.StatusBadGateway)
		return
	}

	// Copy headers from the original request
	for key, vals := range r.Header {
		for _, v := range vals {
			req.Header.Add(key, v)
		}
	}

	resp, err := h.client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", target.String()).Msg("Upstream request failed")
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				log.Warn().Err(writeErr).Msg("Failed to write response to client")
				return
			}
		}
		if readErr != nil {
			break
		}
	}
}
