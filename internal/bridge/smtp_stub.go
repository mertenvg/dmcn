package bridge

import (
	"context"
	"sync"
)

// DeliveredMessage records a message captured by StubSMTPDeliverer.
type DeliveredMessage struct {
	From    string
	To      string
	Subject string
	Body    string
}

// StubSMTPDeliverer is an SMTPDeliverer that captures messages in memory
// for test assertions instead of delivering via SMTP.
type StubSMTPDeliverer struct {
	mu       sync.Mutex
	Messages []DeliveredMessage
}

// Deliver records the message for later inspection.
func (s *StubSMTPDeliverer) Deliver(_ context.Context, from, to, subject, body string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages = append(s.Messages, DeliveredMessage{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
	})
	return nil
}
