package events

import "context"

// metadataKey is the context key for headers.
type metadataKey struct{}

// ContextWithHeaders attaches headers to ctx.
func ContextWithHeaders(ctx context.Context, headers map[string]string) context.Context {
	if len(headers) == 0 {
		return ctx
	}
	// Make a copy to prevent external modifications
	cp := make(map[string]string, len(headers))
	for k, v := range headers {
		cp[k] = v
	}
	return context.WithValue(ctx, metadataKey{}, cp)
}

// HeadersFrom extracts headers from ctx if present.
func HeadersFrom(ctx context.Context) (map[string]string, bool) {
	v := ctx.Value(metadataKey{})
	if v == nil {
		return nil, false
	}
	headers, ok := v.(map[string]string)
	return headers, ok && len(headers) > 0
}
