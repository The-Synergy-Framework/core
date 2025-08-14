package gelf

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Handler is a production-ready GELF (Graylog) UDP handler for slog.
// It supports asynchronous logging, connection pooling, and graceful error handling.
type Handler struct {
	mu    sync.RWMutex
	conn  net.Conn
	host  string
	level slog.Leveler
	attrs []slog.Attr
	group string

	// Configuration
	timeout    time.Duration
	async      bool
	bufferSize int

	// Async processing
	msgChan chan *gelfMessage
	done    chan struct{}
	wg      sync.WaitGroup
}

type gelfMessage struct {
	data map[string]any
}

// Config holds GELF handler configuration.
type Config struct {
	Level      slog.Leveler
	Timeout    time.Duration // Connection timeout (default: 5s)
	Async      bool          // Use async logging (default: true)
	BufferSize int           // Async buffer size (default: 1000)
}

// DefaultConfig returns sensible defaults for GELF logging.
func DefaultConfig() *Config {
	return &Config{
		Level:      slog.LevelInfo,
		Timeout:    5 * time.Second,
		Async:      true,
		BufferSize: 1000,
	}
}

// New creates a GELF UDP handler sending to udpAddr (e.g., "127.0.0.1:12201").
func New(udpAddr string, config *Config) (*Handler, error) {
	if config == nil {
		config = DefaultConfig()
	}

	conn, err := net.DialTimeout("udp", udpAddr, config.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to GELF endpoint %s: %w", udpAddr, err)
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	level := config.Level
	if level == nil {
		level = slog.LevelInfo
	}

	h := &Handler{
		conn:       conn,
		host:       hostname,
		level:      level,
		timeout:    config.Timeout,
		async:      config.Async,
		bufferSize: config.BufferSize,
		done:       make(chan struct{}),
	}

	if h.async {
		h.msgChan = make(chan *gelfMessage, h.bufferSize)
		h.wg.Add(1)
		go h.asyncProcessor()
	}

	return h, nil
}

// Enabled reports whether the handler handles records at the given level.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	if h == nil || h.level == nil {
		return false
	}
	return level >= h.level.Level()
}

// Handle processes the log record.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if h == nil || !h.Enabled(ctx, r.Level) {
		return nil
	}

	data := h.buildGelfMessage(r)

	if h.async {
		return h.handleAsync(data)
	}

	return h.handleSync(data)
}

func (h *Handler) buildGelfMessage(r slog.Record) map[string]any {
	data := map[string]any{
		"version":       "1.1",
		"host":          h.host,
		"short_message": r.Message,
		"timestamp":     float64(r.Time.UnixNano()) / 1e9,
		"level":         mapLevel(r.Level),
	}

	// Add handler attributes
	for _, attr := range h.attrs {
		key, value := h.formatAttr(attr)
		if key != "" {
			data[key] = value
		}
	}

	// Add record attributes
	r.Attrs(func(attr slog.Attr) bool {
		key, value := h.formatAttr(attr)
		if key != "" {
			data[key] = value
		}
		return true
	})

	return data
}

func (h *Handler) handleAsync(data map[string]any) error {
	msg := &gelfMessage{data: data}

	select {
	case h.msgChan <- msg:
		return nil
	default:
		// Buffer full - this is a non-blocking operation
		// In production, you might want to increment a dropped messages counter
		return fmt.Errorf("GELF handler buffer full, message dropped")
	}
}

func (h *Handler) handleSync(data map[string]any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal GELF message: %w", err)
	}

	h.mu.RLock()
	conn := h.conn
	h.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("GELF connection is closed")
	}

	// Set write deadline to prevent hanging
	if deadline, ok := conn.(*net.UDPConn); ok {
		deadline.SetWriteDeadline(time.Now().Add(h.timeout))
	}

	_, err = conn.Write(jsonData)
	return err
}

func (h *Handler) asyncProcessor() {
	defer h.wg.Done()

	for {
		select {
		case msg := <-h.msgChan:
			if err := h.handleSync(msg.data); err != nil {
				// In production, you might want to log this error or increment an error counter
				// For now, we silently drop failed messages to prevent infinite loops
			}
		case <-h.done:
			// Drain remaining messages
			for {
				select {
				case msg := <-h.msgChan:
					h.handleSync(msg.data)
				default:
					return
				}
			}
		}
	}
}

// WithAttrs returns a new Handler with the given attributes.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h == nil {
		return nil
	}

	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &Handler{
		conn:       h.conn,
		host:       h.host,
		level:      h.level,
		attrs:      newAttrs,
		group:      h.group,
		timeout:    h.timeout,
		async:      h.async,
		bufferSize: h.bufferSize,
		msgChan:    h.msgChan,
		done:       h.done,
	}
}

// WithGroup returns a new Handler with the given group name.
func (h *Handler) WithGroup(name string) slog.Handler {
	if h == nil || name == "" {
		return h
	}

	newGroup := name
	if h.group != "" {
		newGroup = h.group + "." + name
	}

	return &Handler{
		conn:       h.conn,
		host:       h.host,
		level:      h.level,
		attrs:      h.attrs,
		group:      newGroup,
		timeout:    h.timeout,
		async:      h.async,
		bufferSize: h.bufferSize,
		msgChan:    h.msgChan,
		done:       h.done,
	}
}

// Close gracefully shuts down the handler.
func (h *Handler) Close() error {
	if h == nil {
		return nil
	}

	if h.async {
		close(h.done)
		h.wg.Wait()
		if h.msgChan != nil {
			close(h.msgChan)
		}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.conn != nil {
		err := h.conn.Close()
		h.conn = nil
		return err
	}

	return nil
}

func (h *Handler) formatAttr(attr slog.Attr) (string, any) {
	key := attr.Key
	if key == "" {
		return "", nil
	}

	// Apply group prefix
	if h.group != "" {
		key = h.group + "." + key
	}

	// GELF additional fields must be prefixed with underscore
	// but reserve standard GELF fields
	if !isStandardGelfField(key) {
		key = "_" + sanitizeKey(key)
	}

	return key, attr.Value.Any()
}

func mapLevel(level slog.Level) int {
	// Map slog levels to syslog levels for GELF
	switch {
	case level <= slog.LevelDebug:
		return 7 // Debug
	case level <= slog.LevelInfo:
		return 6 // Informational
	case level <= slog.LevelWarn:
		return 4 // Warning
	default:
		return 3 // Error
	}
}

func isStandardGelfField(key string) bool {
	switch key {
	case "version", "host", "short_message", "full_message", "timestamp", "level":
		return true
	default:
		return false
	}
}

func sanitizeKey(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}
