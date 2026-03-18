package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/server"
	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/store"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	ss := store.NewSessionStore(time.Hour)
	token, _ := ss.Create("alice@dmcn.me")

	var gotAddr string
	handler := server.AuthMiddleware(ss)(func(w http.ResponseWriter, r *http.Request) {
		gotAddr = store.AddressFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if gotAddr != "alice@dmcn.me" {
		t.Fatalf("expected alice@dmcn.me in context, got %q", gotAddr)
	}
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	ss := store.NewSessionStore(time.Hour)
	handler := server.AuthMiddleware(ss)(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	ss := store.NewSessionStore(time.Hour)
	handler := server.AuthMiddleware(ss)(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer badtoken")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_ExpiredSession(t *testing.T) {
	ss := store.NewSessionStore(time.Millisecond)
	token, _ := ss.Create("alice@dmcn.me")
	time.Sleep(5 * time.Millisecond)

	handler := server.AuthMiddleware(ss)(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestCORSMiddleware_DevMode(t *testing.T) {
	handler := server.CORSMiddleware(true, "dmcn.me")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if origin := rr.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Fatalf("expected * origin in dev mode, got %q", origin)
	}
}

func TestCORSMiddleware_ProdMode(t *testing.T) {
	handler := server.CORSMiddleware(false, "dmcn.me")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if origin := rr.Header().Get("Access-Control-Allow-Origin"); origin != "https://dmcn.me" {
		t.Fatalf("expected https://dmcn.me, got %q", origin)
	}
}

func TestCORSMiddleware_Options(t *testing.T) {
	handler := server.CORSMiddleware(true, "dmcn.me")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called for OPTIONS")
	}))

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
	if methods := rr.Header().Get("Access-Control-Allow-Methods"); methods == "" {
		t.Fatal("expected Allow-Methods header")
	}
}

func TestCSPMiddleware(t *testing.T) {
	handler := server.CSPMiddleware("dmcn.me")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	csp := rr.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "default-src 'self'") {
		t.Fatalf("expected CSP header, got %q", csp)
	}
	if !strings.Contains(csp, "wss://dmcn.me") {
		t.Fatalf("expected wss connect-src, got %q", csp)
	}
}

func TestRateLimitMiddleware_UnderLimit(t *testing.T) {
	limiter := server.RateLimitMiddleware(5)
	handler := limiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "1.2.3.4:5678"
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, rr.Code)
		}
	}
}

func TestRateLimitMiddleware_OverLimit(t *testing.T) {
	limiter := server.RateLimitMiddleware(3)
	handler := limiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "1.2.3.4:5678"
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	// 4th request should be rate limited.
	req := httptest.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "1.2.3.4:5678"
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rr.Code)
	}
}
