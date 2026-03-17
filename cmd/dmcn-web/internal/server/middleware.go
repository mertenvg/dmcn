package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/store"
)

// contextKey is a private type for context value keys in this package.
type contextKey string

// ContextKeyAddress is the context key under which the authenticated user's
// address is stored after successful session validation.
const ContextKeyAddress contextKey = "address"

// AddressFromContext extracts the authenticated address from the request
// context. Returns an empty string if no address is present.
func AddressFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ContextKeyAddress).(string)
	return v
}

// AuthMiddleware returns a middleware that validates the Bearer token from
// the Authorization header against the session store and injects the
// associated address into the request context.
func AuthMiddleware(sessions *store.SessionStore) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" || token == authHeader {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			address, err := sessions.Validate(token)
			if err != nil {
				http.Error(w, `{"error":"invalid or expired session"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyAddress, address)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// CORSMiddleware returns a middleware that sets CORS headers. In dev mode
// all origins are allowed; in production only the configured domain is
// permitted.
func CORSMiddleware(devMode bool, domain string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := "*"
			if !devMode && domain != "" {
				origin = fmt.Sprintf("https://%s", domain)
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CSPMiddleware returns a middleware that sets Content-Security-Policy
// headers appropriate for the configured domain.
func CSPMiddleware(domain string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			csp := "default-src 'self'"
			if domain != "" {
				csp += fmt.Sprintf("; connect-src 'self' wss://%s", domain)
			}
			w.Header().Set("Content-Security-Policy", csp)
			next.ServeHTTP(w, r)
		})
	}
}

// rateLimitEntry tracks request timestamps for a single IP.
type rateLimitEntry struct {
	mu        sync.Mutex
	timestamps []time.Time
}

// RateLimitMiddleware returns a middleware that limits requests per IP
// address to the specified number per minute using a sliding window.
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	var clients sync.Map // IP string → *rateLimitEntry

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			val, _ := clients.LoadOrStore(ip, &rateLimitEntry{})
			entry := val.(*rateLimitEntry)

			now := time.Now()
			windowStart := now.Add(-time.Minute)

			entry.mu.Lock()

			// Prune timestamps outside the sliding window.
			valid := entry.timestamps[:0]
			for _, ts := range entry.timestamps {
				if ts.After(windowStart) {
					valid = append(valid, ts)
				}
			}
			entry.timestamps = valid

			if len(entry.timestamps) >= requestsPerMinute {
				entry.mu.Unlock()
				http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}

			entry.timestamps = append(entry.timestamps, now)
			entry.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
