// Package server provides the HTTP server, routing, and TLS configuration
// for the DMCN web client backend.
package server

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/mertenvg/logr/v2"
	"golang.org/x/crypto/acme/autocert"

	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/api"
)

// Config holds HTTP server configuration.
type Config struct {
	ListenAddr string
	Domain     string
	TLSCert    string
	TLSKey     string
	DevMode    bool
	DataDir    string
}

// Server wraps the standard library HTTP server with DMCN routing.
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	log        logr.Logger
	devMode    bool
	domain     string
}

// New creates a Server with the given configuration and logger.
func New(cfg Config, log logr.Logger) *Server {
	mux := http.NewServeMux()
	handler := CSPMiddleware(cfg.Domain)(CORSMiddleware(cfg.DevMode, cfg.Domain)(mux))
	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.ListenAddr,
			Handler: handler,
		},
		mux:     mux,
		log:     log,
		devMode: cfg.DevMode,
		domain:  cfg.Domain,
	}
}

// RegisterAPI wires API handlers, the WebSocket endpoint, and the embedded
// frontend into the server's multiplexer. The authMiddleware function wraps
// handlers that require an authenticated session.
func (s *Server) RegisterAPI(
	auth *api.AuthHandler,
	msg *api.MessageHandler,
	ident *api.IdentityHandler,
	contacts *api.ContactHandler,
	wsHandler http.HandlerFunc,
	authMiddleware func(http.HandlerFunc) http.HandlerFunc,
	frontendFS fs.FS,
) {
	// Public endpoints with rate limiting.
	rateLimiter := RateLimitMiddleware(20)
	s.mux.Handle("POST /api/v1/register", rateLimiter(http.HandlerFunc(auth.HandleRegister)))
	s.mux.Handle("POST /api/v1/login", rateLimiter(http.HandlerFunc(auth.HandleLogin)))
	s.mux.Handle("POST /api/v1/login/verify", rateLimiter(http.HandlerFunc(auth.HandleLoginVerify)))
	s.mux.Handle("GET /api/v1/relay-hints", rateLimiter(http.HandlerFunc(ident.HandleRelayHints)))

	// Authenticated endpoints.
	s.mux.HandleFunc("POST /api/v1/logout", authMiddleware(auth.HandleLogout))
	s.mux.HandleFunc("POST /api/v1/messages/send", authMiddleware(msg.HandleSend))
	s.mux.HandleFunc("GET /api/v1/messages", authMiddleware(msg.HandleList))
	s.mux.HandleFunc("POST /api/v1/messages/ack", authMiddleware(msg.HandleAck))
	s.mux.HandleFunc("GET /api/v1/identity/lookup", authMiddleware(ident.HandleLookup))
	s.mux.HandleFunc("PUT /api/v1/user/payload", authMiddleware(contacts.HandleUpdatePayload))

	// WebSocket endpoint.
	s.mux.HandleFunc("GET /ws", wsHandler)

	// Embedded frontend — serve static files, fall back to index.html for
	// client-side routing.
	if frontendFS != nil {
		fileServer := http.FileServerFS(frontendFS)
		s.mux.Handle("/", fileServer)
	}
}

// Start begins listening. If certFile and keyFile are non-empty it starts
// a TLS server; otherwise it starts a plain HTTP server.
func (s *Server) Start(certFile, keyFile string) error {
	if certFile != "" && keyFile != "" {
		s.log.Info("starting HTTPS server", logr.M("addr", s.httpServer.Addr))
		return s.httpServer.ListenAndServeTLS(certFile, keyFile)
	}
	s.log.Info("starting HTTP server", logr.M("addr", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

// StartAutocert begins listening with automatic TLS certificates from
// Let's Encrypt via the ACME protocol.
func (s *Server) StartAutocert(domain, cacheDir string) error {
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache(cacheDir),
	}
	s.httpServer.TLSConfig = m.TLSConfig()
	s.log.Info("starting HTTPS server with autocert", logr.M("addr", s.httpServer.Addr), logr.M("domain", domain))
	return s.httpServer.ListenAndServeTLS("", "")
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
