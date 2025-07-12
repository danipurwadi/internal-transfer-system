package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Respond converts a Go value to JSON and sends it to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	setStatusCode(ctx, statusCode)

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// If there is nothing to marshal then set status code and return.
	if data == nil {
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("web.respond: marshal: %w", err)
	}

	if _, err := w.Write(jsonData); err != nil {
		return fmt.Errorf("web.respond: write: %w", err)
	}

	return nil
}
