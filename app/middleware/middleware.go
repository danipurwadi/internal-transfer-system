package middleware

import "context"

// Handler represents the handler function that needs to be called.
type Handler func(context.Context) error
