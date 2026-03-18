package store

import "context"

// contextKey is a private type for context value keys, preventing collisions
// with plain string keys used by other packages.
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
