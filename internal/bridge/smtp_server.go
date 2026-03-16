package bridge

import (
	"bytes"
	"context"
	"io"
	"net"
	"sync"

	"github.com/emersion/go-smtp"
	"github.com/mertenvg/logr/v2"
)

// SMTPServer wraps go-smtp to receive inbound legacy email and pass it
// to the InboundHandler.
type SMTPServer struct {
	server  *smtp.Server
	handler *InboundHandler
	ctx     context.Context
	log     logr.Logger
	mu      sync.Mutex
	started bool
}

// NewSMTPServer creates a new SMTP server that delivers inbound messages
// to the provided InboundHandler.
func NewSMTPServer(ctx context.Context, addr string, handler *InboundHandler, domain string, log logr.Logger) *SMTPServer {
	s := &SMTPServer{
		handler: handler,
		ctx:     ctx,
		log:     log,
	}

	be := &smtpBackend{handler: handler, ctx: ctx, log: log}
	srv := smtp.NewServer(be)
	srv.Addr = addr
	srv.Domain = domain
	srv.AllowInsecureAuth = true // PoC: no TLS requirement
	srv.MaxMessageBytes = 4 * 1024 * 1024

	s.server = srv
	return s
}

// Start begins listening for SMTP connections in a goroutine.
func (s *SMTPServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return nil
	}

	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	s.server.Addr = ln.Addr().String()
	s.started = true
	s.log.Infof("SMTP server listening on %s", s.server.Addr)

	go func() {
		if err := s.server.Serve(ln); err != nil {
			s.log.Debugf("SMTP server stopped: %v", err)
		}
	}()

	return nil
}

// Addr returns the listen address. Only valid after Start().
func (s *SMTPServer) Addr() string {
	return s.server.Addr
}

// Stop gracefully shuts down the SMTP server.
func (s *SMTPServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.started {
		return nil
	}
	s.started = false
	s.log.Info("SMTP server stopped")
	return s.server.Close()
}

// smtpBackend implements smtp.Backend.
type smtpBackend struct {
	handler *InboundHandler
	ctx     context.Context
	log     logr.Logger
}

func (b *smtpBackend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &smtpSession{
		handler:  b.handler,
		ctx:      b.ctx,
		log:      b.log,
		remoteIP: c.Conn().RemoteAddr().String(),
	}, nil
}

// smtpSession implements smtp.Session, collecting MAIL FROM, RCPT TO,
// and DATA before passing to the InboundHandler.
type smtpSession struct {
	handler  *InboundHandler
	ctx      context.Context
	log      logr.Logger
	from     string
	to       string
	remoteIP string
}

func (s *smtpSession) Mail(from string, _ *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *smtpSession) Rcpt(to string, _ *smtp.RcptOptions) error {
	s.to = to
	return nil
}

func (s *smtpSession) Data(r io.Reader) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return err
	}

	if err := s.handler.HandleMessage(
		s.ctx, s.remoteIP, s.from, s.to, buf.Bytes(),
	); err != nil {
		s.log.Warnf("inbound message handling failed: %v", err)
		return err
	}

	return nil
}

func (s *smtpSession) Reset() {
	s.from = ""
	s.to = ""
}

func (s *smtpSession) Logout() error {
	return nil
}
