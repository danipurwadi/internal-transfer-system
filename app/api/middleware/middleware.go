// Package middleware contains shared middleware for the application.
package middleware

import "context"

// Handler represents the handler function that needs to be called.
type Handler func(context.Context) error
